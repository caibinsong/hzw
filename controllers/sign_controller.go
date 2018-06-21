package controllers

import (
	"../db"
	"log"
	"net/http"
	"sync"
)

type SignController struct {
	BaseConttroller
}

var sign_mut sync.Mutex

//查询是否已经签到
func (this *SignController) IsSign(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	isSign, err := db.IsSign(par_map["userid"], par_map["signtype"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, isSign, nil)
}

//签到
func (this *SignController) Sign(rw http.ResponseWriter, req *http.Request) {
	sign_mut.Lock()
	defer sign_mut.Unlock()

	par_map := this.GetParameter(req)

	//判断是否可以签到
	isSign, err := db.IsSign(par_map["userid"], par_map["signtype"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}

	if isSign {
		this.ResultString(rw, nil, this.ACTIONERR())
	} else {
		//签到
		err := db.Sign(par_map["userid"], par_map["signtype"])
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
		//获取礼物查询签到次数
		rst_gift, err := db.SignGift(par_map["userid"], par_map["signtype"])
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
		this.ResultString(rw, rst_gift["gift"], nil)
	}
}
