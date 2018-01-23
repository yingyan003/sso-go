package sysinit

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

var Cache cache.Cache

func CacheInit() {
	var err error
	adapter := beego.AppConfig.String("cache.adapter")
	conn := beego.AppConfig.String("cache.conn")
	key := beego.AppConfig.String("cache.key")

	Cache, err = cache.NewCache(adapter, `{"conn":"`+conn+`",`+`"key"`+`:"`+key+`"}`)
	if err != nil {
		logs.Error("初始化cache配置失败，err=%v", err)
	}
}
