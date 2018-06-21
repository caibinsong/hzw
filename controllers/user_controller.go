package controllers

import (
	"../db"
	"log"
	"net/http"
)

type UserController struct {
	BaseConttroller
}

/*
用户注册
参数:user pwd（用户名密码）
返回：success 或者 err:开头的内容
*/
func (this *UserController) AddNewUser(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	err := db.AddNewUser(par_map["user"], par_map["pwd"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, err)
		return
	}
	this.ResultString(rw, nil, nil)
}

/*
用户登录
参数:user pwd
返回：success 或者 err:开头的内容
*/
func (this *UserController) Login(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	userinfo, err := db.UserLogin(par_map["user"], par_map["pwd"])
	if err != nil {
		this.ResultString(rw, "", err)
		return
	}
	userinfo["pwd"] = ""
	this.ResultString(rw, userinfo, err)
}

/*
用户退出
*/
func (this *UserController) Logout(rw http.ResponseWriter, req *http.Request) {
	this.ResultString(rw, nil, nil)
}

func (this *UserController) InitRole(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	isInitRole := db.IsInitRole(par_map["userid"], par_map["user"])
	if isInitRole {
		log.Println(par_map, "角色初始化")
		for _, one_her := range db.Hero_Base_Info {
			if one_her["name"] == par_map["hero_name"] {
				db.AddHero(par_map["userid"], one_her)
				break
			}
		}
		err := db.ChangeName(par_map["userid"], par_map["user"], par_map["newname"])
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, err)
			return
		}
		this.ResultString(rw, nil, nil)
		return
	} else {
		log.Println(par_map, "不是初始化角色")
		this.ResultString(rw, nil, this.ACTIONERR())
	}
}

//修改角色name
func (this *UserController) ChangeName(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	log.Println(par_map)
	err := db.ChangeName(par_map["userid"], par_map["user"], par_map["newname"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, err)
		return
	}
	this.ResultString(rw, nil, nil)
}
