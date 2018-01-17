package models

import (
	"github.com/pborman/uuid"
	"sso/common/util"
)

const (
	//session有效期（24*60*60=86400）
	DEFAULT_EXPIRATION = 3600
)

type Session struct {
	Id         string `json:"sessionId"`
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
	Createtime string `json:"createtime"`
	//过期时间 todo 这是自定义的，应该不是设置session过期的地方，有待考量
	Expiration int    `json:"expiration"`
}

func NewSession(userId int64, username string) *Session {
	return &Session{
		Id:         uuid.New(),
		UserId:     userId,
		Username:   username,
		Createtime: util.GetCurrentTime(),
		Expiration: DEFAULT_EXPIRATION,
	}
}
