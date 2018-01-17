package errors

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type Status struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	//操作成功
	OK = "OK"
	UNMARSHAL_REQBODY_ERR = "unmarshal_reqbody_err"
	//数据库增删改查错误
	DB_CRUD_ERR = "db_crud_err"
	NEW_USER_ERR = "new_user_err"

	//jwt
	TOKEN_ERR = "token_err"
)

func NewStatus(code, msg string) *Status {
	return &Status{Code: code, Message: msg}
}

//把对象转换为json
func (s *Status) ToString() string {
	data, err := json.Marshal(s)
	if err != nil {
		logs.Error("status marshal tostring error:", err)
		return ""
	}
	return string(data)
}

//把对象转换成[]byte类型，用于匹配http返回的类型
func (s *Status) ToBytes() []byte {
	data, err := json.Marshal(s)
	if err != nil {
		logs.Error("status marshal tostring error:", err)
		return []byte("")
	}
	return data
}
