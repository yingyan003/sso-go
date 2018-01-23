package main

import (
	//必须引用，否则会报错 cache: unknown adapter name "redis" (forgot to import?)
	_ "github.com/astaxie/beego/cache/redis"
	_ "sso/common/sysinit"
	_ "sso/routers"
	"github.com/astaxie/beego"
	"fmt"
	"time"
	"sso/common/sysinit"
	"github.com/garyburd/redigo/redis"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()
}

func DBTest(){
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

	//ca,err:=cache.NewCache("redis", `{"conn":"10.151.30.50:6379","key":"super-zxy"}`)
	//////ca,err:=cache.NewCache("redis", `{"conn":"127.0.0.1:6379","key":"zxy"}`)
	//if err!=nil{
	//	fmt.Print(err)
	//}
}

func RedisTest(){

	//bm, err := cache.NewCache("redis", `{"conn":"10.151.30.50:6379"}`)
	//if err!=nil{
	//	fmt.Println(err)
	//}
	sysinit.Cache.Put("user","zxy",time.Second)
	sysinit.Cache.Put("user","zxy",time.Hour)
	time.Sleep(time.Second*10)
	user,_:=redis.String(sysinit.Cache.Get("user"),nil)
	fmt.Println("change time:",user)

	/*sysinit.Cache.Put("aa","zz",time.Hour)
	v,err:=redis.String(sysinit.Cache.Get("aa"), nil)
	fmt.Println(v)
	if err = sysinit.Cache.Put("q", 1, time.Hour); err != nil {
		fmt.Println(err)
	}
	v2,_:=redis.Int(sysinit.Cache.Get("q"),nil)
	fmt.Println(v2)
	if sysinit.Cache.IsExist("q"){
		sysinit.Cache.Delete("q")
	}
	if !sysinit.Cache.IsExist("q"){
		fmt.Println("q is not exist")
	}
	sysinit.Cache.Put("c","cc",time.Second)
	v3,err:=redis.String(sysinit.Cache.Get("c"),nil)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println("c",v3)
	time.Sleep(time.Second)
	if !sysinit.Cache.IsExist("c"){
		fmt.Println("c is not exist")
	}*/
}