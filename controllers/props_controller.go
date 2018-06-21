package controllers

import (
	"../db"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type PropsController struct {
	BaseConttroller
}

const (
	ADD_HP   = "+hp"
	ADD_ATK  = "+atk"
	ADD_DEF  = "+def"
	ADD_WILL = "+will"
)

var (
	qualityMap map[string]int = map[string]int{"A": 10, "B": 4, "C": 2}
)

//添加道具
func (this *PropsController) AddUserProps(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	props := db.GetPropsInfoByName(par_map["name"], par_map["quality"])
	if props == nil {
		log.Println(par_map["name"], " not find")
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	err := db.AddUserProps(par_map["userid"], props, par_map["lv"], par_map["num"])
	if err != nil {
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, nil, nil)
}

func (this *PropsController) UseProps(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	if par_map["name"] == "潜力珠" {
		userheroinfo, err := db.QueryMyHeroInfoById(par_map["userid"], par_map["userheroid"])
		if err != nil {
			log.Println(err)
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
		if userheroinfo["pot"] == "0" {
			this.ResultString(rw, "潜力为0", nil)
			return
		}
		act, err := this.qlzHero(userheroinfo, par_map["quality"], par_map["num"])
		if err != nil {
			log.Println(err)
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		} else {
			md5, err := db.UseQLZ(par_map["userid"], par_map, par_map["num"], act)
			if err != nil {
				log.Println(err)
				this.ResultString(rw, nil, this.ACTIONERR())
				return
			}
			this.ResultString(rw, fmt.Sprintf("%s md5:%s", act, md5), nil)
			return
		}
	} else {
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, nil, nil)
}

func (this *PropsController) qlzHero(heroinfo map[string]string, quality, num string) (string, error) {
	_num, err := strconv.Atoi(num)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	base, ok := qualityMap[quality]
	if !ok {
		log.Println("quality not in (A,B,C)")
		return "", err
	}
	up_num := _num * base
	pot, err := strconv.Atoi(heroinfo["pot"])
	if err != nil {
		log.Println("heroinfo[pot] string to int err:", err.Error())
		return "", err
	}
	if pot < up_num {
		up_num = pot
	}
	switch this.Rand(4) {
	case 0:
		return this.fmt_act(heroinfo["id"], ADD_HP, up_num), nil
	case 1:
		return this.fmt_act(heroinfo["id"], ADD_ATK, up_num), nil
	case 2:
		return this.fmt_act(heroinfo["id"], ADD_DEF, up_num), nil
	case 3:
		return this.fmt_act(heroinfo["id"], ADD_WILL, up_num), nil
	default:
		return this.fmt_act(heroinfo["id"], ADD_HP, up_num), nil
	}
}

func (this *PropsController) fmt_act(id, types string, num int) string {
	return fmt.Sprintf("heroid:%s %s:%d", id, types, num)
}
