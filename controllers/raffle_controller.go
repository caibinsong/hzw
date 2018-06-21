package controllers

import (
	"../db"
	"../models"
	"log"
	"net/http"
	"strconv"
)

//抽奖
type RaffleController struct {
	BaseConttroller
}

func (this *RaffleController) raffle_base(req *http.Request) (map[string]string, error) {
	par_map := this.GetParameter(req)
	userid_int, err := strconv.Atoi(par_map["userid"])
	if err != nil {
		return par_map, err
	}
	used := par_map["used"]
	used_int, err := strconv.Atoi(used)
	if err != nil {
		return par_map, err
	}
	//花费钻石
	err = db.AddMoney(userid_int, 0, used_int, 0)
	if err != nil {
		log.Println(err.Error())
		return par_map, err
	}
	return par_map, nil
}

//抽英雄
func (this *RaffleController) DoHero(rw http.ResponseWriter, req *http.Request) {
	par_map, err := this.raffle_base(req)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	//随机出英雄
	hero := models.RandGet(par_map["quality"], db.Hero_Base_Info)
	//判断是否已经有了
	err = db.AddHero(par_map["userid"], hero)

	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, hero, nil)
}

//抽装备
func (this *RaffleController) DoEqpt(rw http.ResponseWriter, req *http.Request) {
	par_map, err := this.raffle_base(req)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	//随机出装备
	eqpt := models.RandGet(par_map["quality"], db.Eqpt_Base_Info)

	err = db.AddEqpt(par_map["userid"], eqpt)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, eqpt, nil)
}

//抽技能
func (this *RaffleController) DoSkill(rw http.ResponseWriter, req *http.Request) {
	par_map, err := this.raffle_base(req)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	//随机出装备
	skill := models.RandGet(par_map["quality"], db.Skill_Base_Info)

	err = db.AddSkill(par_map["userid"], skill)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, skill, nil)
}
