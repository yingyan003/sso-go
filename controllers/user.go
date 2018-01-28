package controllers

import (
	"sso/models"
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"sso/common/util/encrypt"
	"strconv"
)

type UserController struct {
	BaseController
}

//创建用户
func (u *UserController) Post() {
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
	param := u.Ctx.Input.Param(":idOrList")
	if param == "list" {
		u.GetUserList()
		return
	} else {
		u.GetUserOne(param)
		return
	}
}

func (u *UserController) GetUserList() {
	users, status := models.QueryUserList()
	if status != nil {
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	logs.Info("查询用户列表成功: 总记录数=%d", len(users))
	status = errors.NewStatus(errors.OK, "查询用户列表成功")
	status.Data = users
	u.Ctx.Output.Body(status.ToBytes())
}

func (u *UserController) GetUserOne(uid string) {
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logs.Error("请求的uid格式错误:uid=%v", uid)
		status := errors.NewStatus(errors.REQUEST_PARAM_ERR, "请求的uid格式错误")
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	user := new(models.User)
	user, status := models.QueryUserById(id)
	if status != nil {
		logs.Info("用户不存在,id=%d", id)
		u.Ctx.Output.Body(status.ToBytes())
		return
	}

	status = errors.NewStatus(errors.OK, "用户存在")
	status.Data = user
	u.Ctx.Output.Body(status.ToBytes())
}
