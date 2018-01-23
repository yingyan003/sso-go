package controllers

import (
	_ "github.com/astaxie/beego/session/redis"
	"sso/models"
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"fmt"
	"sso/common/util/encrypt"
	"github.com/pborman/uuid"
	"sso/common/sysinit"
	"time"
	"sso/common/util/jwtUtil"
)

const (
	//登录状态有效期  24*60*60s=86400
	DEFAULT_EXPIRATION = 4 * time.Hour
)

type AccessController struct {
	BaseController
}

func (a *AccessController) Auth() {
	status := errors.NewStatus(errors.OK, "认证成功，用户已登录")
	status.Data = a.Ctx.Input.Param("userId")
	a.Ctx.Output.Body(status.ToBytes())
}

func (a *AccessController) Login() {
	user := new(models.User)
	status := a.getReqBody(user)
	if status != nil {
		logs.Error("登录失败，username=%s, msg=%s", user.Username, status.Message)
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	user.Password = encrypt.NewEncryption(user.Password).String()
	//检查用户名密码是否正确
	user, status = models.QueryUserByNameAndPwd(user.Username, user.Password)
	if status != nil {
		logs.Error("登录失败，用户名或密码错误")
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	redisKey, err := a.saveStatus(user.Id)
	if err != nil {
		logs.Error("保存登录状态到redis失败，error=%v", err)
		status := errors.NewStatus(errors.REDIS_ERR, "登录状态保存到redis失败")
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	token, status := createToken(redisKey, user)
	if status != nil {
		logs.Error("登录失败，创建token失败")
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	logs.Info("登录成功，username=%s", user.Username)
	status = errors.NewStatus(errors.OK, "登录成功")
	status.Data = token
	a.Ctx.Output.Body(status.ToBytes())
}

func (a *AccessController) Logout() {
	redisKey := a.Ctx.Input.Param("redisKey")
	if !IsLoginStatus(redisKey) {
		logs.Info("Logout:用户登录状态在redis中已失效")
		status := errors.NewStatus(errors.TOKEN_ERR, "token已失效")
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	err := sysinit.Cache.Delete(redisKey)
	if err != nil {
		logs.Error("删除Redis的登录状态失败，error=%v", err)
		status := errors.NewStatus(errors.REDIS_ERR, "登录状态删除失败")
		a.Ctx.Output.Body(status.ToBytes())
		return
	}

	status := errors.NewStatus(errors.OK, "用户退出成功")
	a.Ctx.Output.Body(status.ToBytes())
}

//将用户登录状态保存到Redis中
func (a *AccessController) saveStatus(uid int64) (string, error) {
	key := uuid.New()
	return key, sysinit.Cache.Put(key, uid, DEFAULT_EXPIRATION)
}

func createToken(redisKey string, user *models.User) (string, *errors.Status) {
	tokenFirst, status := jwtUtil.GetEStoken(redisKey)
	if status != nil {
		return "", status
	}
	tokenSecond, status := jwtUtil.GetHStoken(tokenFirst, user)
	if status != nil {
		return "", status
	}
	return tokenSecond, status
}

//判断用户是否处于登录状态
func IsLoginStatus(redisKey string) bool {
	if sysinit.Cache.IsExist(redisKey) {
		return true
	}
	return false
}

//session相关的代码暂时保留
func (a *AccessController) createSession(uid int64, username string) int64 {
	a.SetSession(uid, username)
	return uid
}

func (a *AccessController) deleteSession(sessionId int64) {
	s := a.GetSession(sessionId)
	if s != nil {
		fmt.Println("deltession: session=%v", s)
		a.DelSession(sessionId)
	} else {
		fmt.Println("a.DelSession(sessionId): session not exist")
	}
}
