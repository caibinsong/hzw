package controllers

import (
	"../db"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MoneyController struct {
	BaseConttroller
}

//金币销费 传入userid 和使用的金币数 包括金币钻石和RMB的消费
func (this *MoneyController) Used(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	money := par_map["money"]
	userid_int, err := strconv.Atoi(par_map["userid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	arr_m := strings.Split(money, "!")
	if len(arr_m) != 3 {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	gold, err := strconv.Atoi(arr_m[0])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	diamond, err := strconv.Atoi(arr_m[1])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	rmb, err := strconv.Atoi(arr_m[2])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	err = db.AddMoney(userid_int, gold, diamond, rmb)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, nil, nil)
}
