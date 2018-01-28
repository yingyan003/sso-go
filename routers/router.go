package routers

import (
	"sso/controllers"
	"github.com/astaxie/beego"
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

	beego.InsertFilter("/*", beego.BeforeRouter, Filter)

	beego.Router("/v1/auth", &controllers.AccessController{}, "get:Auth")
	beego.Router("/v1/login", &controllers.AccessController{}, "post:Login")
	beego.Router("/v1/logout", &controllers.AccessController{}, "get:Logout")
	beego.Router("/v1/user/?:idOrList", &controllers.UserController{})
}
