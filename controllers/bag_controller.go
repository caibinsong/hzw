package controllers

import (
	"../db"
	"log"
	"net/http"
)

type BagController struct {
	BaseConttroller
}

func (this *BagController) Bag(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	var rst []map[string]string
	var err error
	if par_map["types"] == "hero" {
		rst, err = db.BagHeroByUserid(par_map["userid"])
	} else if par_map["types"] == "herochip" {
		rst, err = db.BagHeroChipByUserid(par_map["userid"])
	} else if par_map["types"] == "eqpt" {
		rst, err = db.BagEqptByUserid(par_map["userid"], par_map["types2"])
	} else if par_map["types"] == "skill" {
		rst, err = db.BagSkillByUserid(par_map["userid"])
	} else {
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, rst, nil)
}

//查看自己背包里的英雄详细信息
func (this *BagController) MyHeroInfo(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	rst, err := db.QueryMyHeroInfoById(par_map["userid"], par_map["userheroid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, rst, nil)
}
func (this *BagController) MyHeroChipInfo(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	rst, err := db.QueryMyHeroChipInfoById(par_map["userid"], par_map["userherochipid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, rst, nil)
}

//查看自己背包里的装备详细信息
func (this *BagController) MyEqptInfo(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	rst, err := db.QueryEqptInfoById(par_map["userid"], par_map["usereqptid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, rst, nil)
}

//查看自己背包里的技能详细信息
func (this *BagController) MySkillInfo(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	rst, err := db.QuerySkillInfoById(par_map["userid"], par_map["userskillid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, rst, nil)
}
