package timeUtil

import "time"

//获取当期时间，时间格式为datetime
func GetCurrentTime()string{
	return time.Now().Format("2006-01-02 15:04:05")
}
