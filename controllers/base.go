package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"sso/common/errors"
	"github.com/astaxie/beego/logs"
)

type BaseController struct {
	beego.Controller
}

//从http请求体中获取用户信息
func (b *BaseController) getReqBody(obj interface{}) *errors.Status {
	if err := json.Unmarshal(b.Ctx.Input.RequestBody, obj); err != nil {
		logs.Error("getReqBody fail,error:",err)
		status := errors.NewStatus(errors.UNMARSHAL_REQBODY_ERR, "解析请求体的json失败")
		return status
	}
	return nil
}
