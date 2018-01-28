package routers

import (
	//session使用Redis做缓存时，必须导入这个import
	//除了 file，memory 和 cookie 的支持，其他都需要引入 beego/session/ 下的包
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"sso/models"
	"sso/controllers"
	"fmt"
	"sso/common/util/jwtUtil"
	"strconv"
	"sso/common/sysinit"
	"github.com/astaxie/beego/plugins/cors"
)

var crossDomainFilter = cors.Allow(&cors.Options{
	AllowAllOrigins:  true,
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
	AllowCredentials: true,
})

//获取用户列表是否直接放行
func Filter(ctx *context.Context) {
	//解决跨域问题
	crossDomainFilter(ctx)
	//登录认证拦截
	auth(ctx)
}

func auth(ctx *context.Context){
	req := ctx.Request
	if req.RequestURI == "/v1/user" && (req.Method == "POST" || req.Method == "GET") {
		fmt.Println(ctx.Request.RequestURI,"放行")
		//ctx.Redirect(302,"/v1/login")
		return
	}

	tokenString := ctx.Input.Header("token")

	//token值为空，其他都需重新登录
	if tokenString == "" {
		status := errors.NewStatus(errors.TOKEN_ERR, "无token")
		ctx.Output.Body(status.ToBytes())
		return
	}
	//首次登录，放行
	//前端登录请求头带token,且默认值为"default"。
	//当token!="default"说明不是首次登录，需验证是否重复登陆
	if tokenString == "default" && req.RequestURI == "/v1/login" {
		fmt.Println("default")
		return
	}

	//解析第一层token
	claimsHS, status := jwtUtil.ParseHStoken(tokenString)
	if status != nil {
		ctx.Output.Body(status.ToBytes())
		return
	}

	//TODO 第二层token采用ES256算法，但XYD没解决
	tokenES := claimsHS["id"].(string)
	//解析第二层token
	redisKey, status := jwtUtil.ParseEStoken(tokenES)
	if status != nil {
		ctx.Output.Body(status.ToBytes())
		return
	}

	if controllers.IsLoginStatus(redisKey) {
		//重复登录
		if req.RequestURI == "/v1/login" {
			status := errors.NewStatus(errors.TOKEN_ERR, "重复登录")
			ctx.Output.Body(status.ToBytes())
			return
		}
		//保存参数，为控制器使用
		uid := claimsHS["userId"].(float64)
		ctx.Input.SetParam("userId", strconv.FormatFloat(uid, 'f', 0, 64))
		ctx.Input.SetParam("redisKey", redisKey)

		//每次访问都重新刷新token失效时间
		sysinit.Cache.Put(redisKey, uid, controllers.DEFAULT_EXPIRATION)
		return
	}

	//登录时携带的token失效，放行
	if req.RequestURI == "/v1/login" {
		return
	}

	logs.Error("认证失败，token已失效")
	status = errors.NewStatus(errors.TOKEN_ERR, "token已失效")
	ctx.Output.Body(status.ToBytes())
}

func IsSessionExist(ctx *context.Context, sessionId string) bool {
	fmt.Println("sessionId", sessionId)
	//s := ctx.Input.CruSession.Get(sessionId)
	//fmt.Println("getsession,s", s)
	//todo 在controller外使用session
	//session失效时，redis会自动删除，Get()会返回nil；当用nil进行断言，则ok=false
	if _, ok := ctx.Input.CruSession.Get(sessionId).(*models.Session); ok {
		return true
	}
	return false
}
