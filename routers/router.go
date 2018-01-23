package routers

import (
	"sso/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)

	//解决浏览器跨域问题
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "token"},
		AllowCredentials: true,
	}))

	beego.InsertFilter("/*", beego.BeforeRouter, AuthFilter)

	beego.Router("/v1/auth", &controllers.AccessController{}, "get:Auth")
	beego.Router("/v1/login", &controllers.AccessController{}, "post:Login")
	beego.Router("/v1/logout", &controllers.AccessController{}, "get:Logout")
	beego.Router("/v1/user", &controllers.UserController{})
}
