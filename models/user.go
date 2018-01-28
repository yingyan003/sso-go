package models

import (
	"github.com/astaxie/beego/logs"
	"sso/common/errors"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id         int64	`json:"id"`
	Username   string	`json:"username"`
	Password   string	`json:"password"`
	Createtime time.Time	`json:"createtime"`
	Updatetime time.Time	`json:"updatetime"`
}

func AddUser(u *User) (int, *errors.Status) {
	if u.Username == "" {
		logs.Error("AddUser:用户名不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户名不能为空")
		return 0, status
	}
	if u.Password == "" {
		logs.Error("AddUser:用户密码不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户密码不能为空")
		return 0, status
	}
	id, err := orm.NewOrm().Insert(u)
	if err != nil {
		logs.Error("AddApp fail, err:", err)
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户已存在")
		return 0, status
	}
	return int(id), nil
}

func DeleteUser(username string) *errors.Status {
	if username == "" {
		logs.Error("DeleteUser:用户名不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户名不能为空")
		return status
	}
	//todo 用户不存在，删除也不返回error
	_, err := orm.NewOrm().QueryTable("user").Filter("username", username).Delete()
	if err != nil {
		logs.Error("delete user fail,err:", err)
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户删除失败")
		return status
	}
	return nil
}

func UpdateUser(u *User) *errors.Status {
	if u.Username == "" {
		logs.Error("UpdateUser:用户名不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户名不能为空")
		return status
	}
	if u.Password == "" {
		logs.Error("UpdateUser:密码不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "密码不能为空")
		return status
	}
	//todo 当username不存在时，并不返回错误
	_, err := orm.NewOrm().QueryTable("user").Filter("username", u.Username).Update(orm.Params{"password": u.Password})
	if err != nil {
		logs.Error("update user fail,err:", err)
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户更新失败")
		return status
	}
	return nil
}

func QueryUserById(id int64) (*User, *errors.Status) {
	if id == 0 {
		logs.Error("QueryUserById:用户id不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户id不能为空")
		return nil, status
	}
	u := new(User)
	u.Id=id
	err := orm.NewOrm().Read(u)
	//err==orm.ErrNoRows说明找不到该记录
	if err != nil {
		logs.Info("QueryUserById fail,err:", err)
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户不存在")
		return nil, status
	}
	return u, nil
}

func QueryUserByName(username string) (*User, *errors.Status) {
	if username == "" {
		logs.Error("QueryUserByName:用户名不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户名不能为空")
		return nil, status
	}
	u := new(User)
	u.Username=username
	err := orm.NewOrm().Read(u, "username")
	//err==orm.ErrNoRows说明找不到该记录
	if err != nil {
		logs.Info("query user fail,err:", err)
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户不存在")
		return nil, status
	}
	return u, nil
}

func QueryUserByNameAndPwd(username,password string) (*User, *errors.Status) {
	if username == "" {
		logs.Error("QueryUserByNameAndPwd:用户名不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户名不能为空")
		return nil, status
	}
	if password ==""{
		logs.Error("QueryUserByNameAndPwd:密码不能为空")
		status := errors.NewStatus(errors.DB_CRUD_ERR, "密码不能为空")
		return nil, status
	}
	u := new(User)
	u.Username=username
	u.Password=password
	err := orm.NewOrm().Read(u, "username","password")
	//err==orm.ErrNoRows说明找不到该记录
	if err != nil {
		logs.Error("query user fail,err:", err)
		status := errors.NewStatus(errors.DB_CRUD_ERR, "用户名或密码错误")
		return nil, status
	}
	return u, nil
}

func QueryUserList()(users []*User,status *errors.Status){
	//默认查询最大行数为1000
	_,err:=orm.NewOrm().QueryTable("user").All(&users,"id","username")
	if err!=nil{
		logs.Error("QueryUserList fail,err=%v",err)
		status=errors.NewStatus(errors.DB_CRUD_ERR,"查询用户列表失败")
		return nil,status
	}
	return users,nil
}

/*
func QUserByTime(){
	start:="2018-01-11 17:55:07"
	end:="2018-01-11 17:56:57"
	var us []*User
	orm.NewOrm().QueryTable("user").Filter("createtime__gte",start).Filter("createtime__lt",end).All(&us)
	for _,u:=range us{
		fmt.Println(*u)
	}
}*/
