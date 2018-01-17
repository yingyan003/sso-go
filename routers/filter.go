package routers

import (
	//session使用Redis做缓存时，必须导入这个import
	//除了 file，memory 和 cookie 的支持，其他都需要引入 beego/session/ 下的包
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"sso/models"
)

func AuthFilter(ctx *context.Context) {
	req:=ctx.Request
	if req.RequestURI== "/v1/login" || (req.RequestURI=="/v1/user" && req.Method == "POST"){
		//ctx.Redirect(302,"/v1/login")
		return
	}

	//tokenString:=ctx.Input.Header("token")
	//token,err:=jwt.Parse(tokenString,func(token *jwt.Token)(interface{},error){
	//	return []byte(controllers.SIGNED_KEY),nil
	//})
	//fmt.Println("tokenstring",tokenString)
	//
	//token,err:=jwt.Parse(tokenString,nil)
	//if err!=nil{
	//	logs.Error("认证失败，token解析错误，err:",err)
	//	status:=errors.NewStatus(errors.TOKEN_ERR,"token解析失败")
	//	ctx.Output.Body(status.ToBytes())
	//	return
	//}
	//
	//claims,ok:=token.Claims.(jwt.MapClaims)
	//if !ok{
	//	logs.Error("认证登录时，获取jwt claims失败")
	//	status:=errors.NewStatus(errors.TOKEN_ERR,"认证登录时，jwt claimts获取失败")
	//	ctx.Output.Body(status.ToBytes())
	//	return
	//}
	//sessionId:=claims["sessionId"].(string)

	sessionId := ctx.Input.Header("sessionId")
	if IsSessionExist(ctx, sessionId) {
		return
	}

	logs.Error("认证失败，session已失效")
	status := errors.NewStatus(errors.TOKEN_ERR, "认证失败，session已失效")
	ctx.Output.Body(status.ToBytes())
}

func IsSessionExist(ctx *context.Context, sessionId string) bool {
	//todo 在controller外使用session
	//session失效时，redis会自动删除，Get()会返回nil；当用nil进行断言，那ok=false
	if _, ok := ctx.Input.CruSession.Get(sessionId).(*models.Session); ok {
		return true
	}
	return false
}
