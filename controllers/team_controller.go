package controllers

import (
	"../db"
	"log"
	"net/http"
)

//阵容
type TeamController struct {
	BaseConttroller
}

//阵容查询
func (this *TeamController) MyTeam(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	rst := this.TeamByUserid(par_map["userid"])
	if rst == nil {
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, rst, nil)
}

//阵容换人
func (this *TeamController) SetHero(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	err := db.SetHero(par_map["userid"], par_map["userheroid"], par_map["teamindex"])
	if err != nil {
		this.ResultString(rw, nil, err)
		return
	}
	this.ResultString(rw, nil, nil)
}

//换装备 0为下
func (this *TeamController) SetEqpt(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	err := db.SetEqptSkill(par_map["userid"], par_map["teamindex"], par_map["types"], par_map["usereqptid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, err)
		return
	}
	this.ResultString(rw, nil, nil)
}

//换技能 0为下
func (this *TeamController) SetSkill(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	err := db.SetEqptSkill(par_map["userid"], par_map["teamindex"], par_map["types"], par_map["userskillid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, err)
		return
	}
	this.ResultString(rw, nil, nil)
}
