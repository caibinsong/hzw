package controllers

import (
	"../db"
	"net/http"
	"strconv"
)

//体力
type StrengthController struct {
	BaseConttroller
}

//输入用户ID和使用的体力，返回是否成功，使用的体力必须大于可用的体力
func (this *StrengthController) UseStrength(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	userid_int, err := strconv.Atoi(par_map["userid"])
	if err != nil {
		this.ResultString(rw, "", err)
		return
	}
	used_int, err := strconv.Atoi(par_map["used"])
	if err != nil {
		this.ResultString(rw, "", err)
		return
	}
	err = db.UsedStrengthByUserId(userid_int, used_int)
	if err != nil {
		this.ResultString(rw, "", err)
		return
	}
	this.ResultString(rw, nil, nil)
}

//传入userid，查询体力
func (this *StrengthController) QueryStrength(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	userid_int, err := strconv.Atoi(par_map["userid"])
	if err != nil {
		this.ResultString(rw, "", err)
		return
	}
	rst, err := db.QueryStrengthByUserId(userid_int)
	if err != nil {
		this.ResultString(rw, "", err)
		return
	}
	this.ResultString(rw, rst, nil)
}
