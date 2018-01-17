package main

import (
	_ "sso/common/sysinit"
	_ "sso/routers"
	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	//u:=new(models.User)
	//u.Username="y"
	//u.Password="22"
	//models.AddUser(u)



	//models.UpdateUser(u)
	//u,status:=models.QueryUserByName("yy")
	//if status!=nil{
	//	fmt.Println("query error")
	//	return
	//}
	//fmt.Println(*u)
	//models.QUserByTime()


	beego.BConfig.WebConfig.Session.SessionAutoSetCookie=false
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime=0
	beego.Run()
}
