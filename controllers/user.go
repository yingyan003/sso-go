package controllers

import (
	"sso/models"
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"sso/common/util/encrypt"
	"fmt"
)

type UserController struct {
	BaseController
}

//创建用户
func (u *UserController) Post() {
	fmt.Println("enter create user")

	//从请求体中解析用户信息
	user := new(models.User)
	status := u.getReqBody(user)
	if status != nil {
		logs.Error("创建用户失败，username=%s, msg=%s", user.Username, status.Message)
		status := errors.NewStatus(errors.UNMARSHAL_REQBODY_ERR, "解析请求体的json失败")
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	//检查用户是否存在
	us, _ := models.QueryUserByName(user.Username)
	if us != nil {
		logs.Error("用户已存在，username=%s", user.Username)
		status := errors.NewStatus(errors.NEW_USER_ERR, "用户已存在")
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	//用户密码加密
	user.Password = encrypt.NewEncryption(user.Password).String()

	_, status = models.AddUser(user)
	if status != nil {
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	logs.Info("用户创建成功: username=%s, password=%s", user.Username, user.Password)
	status = errors.NewStatus(errors.OK, "用户创建成功")
	u.Ctx.Output.Body(status.ToBytes())
}

//删除用户
func (u *UserController) Delete() {
	user := new(models.User)
	status := u.getReqBody(user)
	if status != nil {
		logs.Error("删除用户失败，username=%s, msg=%s", user.Username, status.Message)
		status := errors.NewStatus(errors.UNMARSHAL_REQBODY_ERR, "解析请求体的json失败")
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	status = models.DeleteUser(user.Username)
	if status != nil {
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	//todo 删除用户时要删除会话

	logs.Info("用户删除成功: username=%s", user.Username)
	status = errors.NewStatus(errors.OK, "用户删除成功")
	u.Ctx.Output.Body(status.ToBytes())
	return
}

//更新用户
func (u *UserController) Put() {
	user := new(models.User)
	status := u.getReqBody(user)
	if status != nil {
		logs.Error("更新用户失败，username=%s, msg=%s", user.Username, status.Message)
		status := errors.NewStatus(errors.UNMARSHAL_REQBODY_ERR, "解析请求体的json失败")
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	user.Password = encrypt.NewEncryption(user.Password).String()
	status = models.UpdateUser(user)
	if status != nil {
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	logs.Info("用户更新成功: username=%s, password=%s", user.Username, user.Password)
	status = errors.NewStatus(errors.OK, "用户更新成功")
	u.Ctx.Output.Body(status.ToBytes())
	return
}

func (u *UserController) Get() {
	users, status := models.QueryUserList()
	if status != nil {
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	logs.Info("查询用户列表成功: 总记录数=%d", len(users))
	status = errors.NewStatus(errors.OK, "查询用户列表成功")
	status.Data = users
	u.Ctx.Output.Body(status.ToBytes())
	return
}

