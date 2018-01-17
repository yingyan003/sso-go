package  sysinit

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"net/url"
	"github.com/astaxie/beego/orm"
	"sso/models"
	"fmt"
	"encoding/gob"
)

const(
	//todo 待确定
	//数据库的空闲连接数
	maxIdle = 30
	//数据库的最大连接
	maxConn = 30

)

func init(){
	fmt.Println("zxy",beego.AppConfig.String("zxy"))

	//获取数据库配置信息
	dbhost:=beego.AppConfig.String("db.host")
	Dbport:=beego.AppConfig.String("db.port")
	dbuser:=beego.AppConfig.String("db.user")
	dbpassword:=beego.AppConfig.String("db.password")
	dbname:=beego.AppConfig.String("db.name")
	timezone:=beego.AppConfig.String("db.timezone")
	if Dbport==""{
		Dbport="3306"
	}

	//组装数据源
	ds:=dbuser+":"+dbpassword+"@tcp("+dbhost+":"+Dbport+")/"+dbname+"?charset=utf8"
	if timezone!=""{
		ds=ds+"&loc="+url.QueryEscape(timezone)
	}

	fmt.Println("ds",ds)

	//todo 解决beego框架本身bug的方式
	//在model里自定义struct，但并不进行数据库操作，故不使用orm.RegisterModel（）注册模型时会报错
	//错误信息如下：gob: name not registered for interface: "*models.Session"
	gob.Register(new(models.Session))

	//注册数据库
	orm.RegisterDataBase("default","mysql",ds)

	//注册模型
	orm.RegisterModel(new(models.User))

	//设置日志级别
	if beego.AppConfig.String("runmode")=="dev"{
		orm.Debug=true
	}
}