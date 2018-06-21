package controllers

import (
	"../db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

//决斗
type DuelController struct {
	BaseConttroller
}

var duel_Mux sync.Mutex
var risk_Mux sync.Mutex

//查询关卡列表 duel/queryadventure
func (this *DuelController) QueryAdventure(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	//先查询自己的排名
	list, err := db.QueryAdventure(par_map["userid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	this.ResultString(rw, list, nil)
}

//挑战关卡
func (this *DuelController) Risk(rw http.ResponseWriter, req *http.Request) {
	risk_Mux.Lock()
	defer risk_Mux.Unlock()
	par_map := this.GetParameter(req)
	//查询出自己阵容
	self_team := this.TeamByUserid(par_map["userid"])
	if self_team == nil {
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	//目前冒险内容未入
	rival_team, err := db.QueryAdventureHero(par_map["adventureindex"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	//结束
	result_list := this.Duel(self_team, rival_team)
	for k, v := range result_list {
		log.Println(k, v)
	}
	//胜利了
	if strings.Index(result_list[len(result_list)-1], "win") >= 0 {
		//修改冒险次数
		err := db.UpdateAdventure(par_map["userid"], par_map["adventureindex"])
		if err != nil {
			log.Println(err.Error())
		}
		//关卡信息
		adventureinfo, err := db.QueryAdventureInfo(par_map["adventureindex"])
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
		//冒险英雄信息
		duelWinHeroInfo, err := db.DuelWinHeroInfo(par_map["userid"])
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
		//经验等级
		this.IsWin(duelWinHeroInfo, adventureinfo)
		//金币
		userid, err := strconv.Atoi(par_map["userid"])
		if err != nil {
			log.Println(err.Error())
		}
		gold, err := strconv.Atoi(adventureinfo["gold"])
		if err != nil {
			log.Println(err.Error())
		}
		err = db.AddMoney(userid, gold, 0, 0)
		if err != nil {
			log.Println(err.Error())
		}
		//是否有掉落
		if strings.Trim(adventureinfo["appear"], " ") != "" {
			appearList := this.WinAppear(par_map["userid"], adventureinfo["appear"])
			if len(appearList) != 0 {
				result_list = append(result_list, "掉落："+strings.Join(appearList, ","))
			}
		}
	}
	this.ResultString(rw, result_list, nil)
}

func (this *DuelController) WinAppear(userid, appear string) []string {
	result := make([]string, 0)
	appearList := strings.Split(appear, "，")
	for _, one_appear := range appearList {
		randnum := this.Rand(100)
		if randnum > 70 {
			eqpt := db.GetEqptInfoByName(one_appear)
			if eqpt != nil {
				err := db.AddEqpt(userid, eqpt)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				result = append(result, one_appear)
				continue
			}
			skill := db.GetSkillInfoByName(one_appear)
			if skill != nil {
				err := db.AddSkill(userid, skill)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				result = append(result, one_appear)
				continue
			}
		}
	}
	return result
}

func (this *DuelController) IsWin(herolist []map[string]string, adventureinfo map[string]string) {
	for k, _ := range herolist {
		nowexp := this.AddIntNum(herolist[k]["nowexp"], adventureinfo["exp"])
		upexp := this.MultiplyIntNum(herolist[k]["lv"], herolist[k]["initexp"])
		if nowexp > upexp {
			nowexp = nowexp - upexp
			err := db.DuelWinHero_UpLv(herolist[k]["hpgrowth"], herolist[k]["atkgrowth"], herolist[k]["defgrowth"],
				herolist[k]["willgrowth"], herolist[k]["potgrowth"],
				strconv.Itoa(nowexp), herolist[k]["id"])
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			err := db.DuelWinHero_UpExp(adventureinfo["exp"], herolist[k]["id"])
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

//查看自己能看到的决斗列表
func (this *DuelController) GetDuelList(rw http.ResponseWriter, req *http.Request) {
	par_map := this.GetParameter(req)
	//先查询自己的排名
	myduel_info_list, err := db.QueryDuelByUserid(par_map["userid"])
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	myIndex := 201
	if len(myduel_info_list) != 0 {
		myIndex, err = strconv.Atoi(myduel_info_list[0]["index"])
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
	}
	//通过自己的排名计算出要显示的排名列表
	result_list, err := db.QueryDuelBetween(1, 5)
	if err != nil {
		log.Println(err.Error())
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	if myIndex > 5 {
		end := myIndex
		start := myIndex - 5
		if start <= 5 {
			start = 6
		}
		append_result_list, err := db.QueryDuelBetween(start, end)
		if err != nil {
			log.Println(err.Error())
			this.ResultString(rw, nil, this.ACTIONERR())
			return
		}
		result_list = append(result_list, append_result_list...)
	}
	log.Println(result_list)
	this.ResultString(rw, result_list, nil)
}

//挑战
func (this *DuelController) Challenge(rw http.ResponseWriter, req *http.Request) {
	duel_Mux.Lock()
	defer duel_Mux.Unlock()
	par_map := this.GetParameter(req)
	//查询出自己阵容
	self_team := this.TeamByUserid(par_map["userid"])
	if self_team == nil {
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	self_index := "0"
	self_info, err := db.QueryDuelByUserid(par_map["userid"])
	if err != nil {
		log.Println("err:", err)
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	if len(self_info) != 0 {
		self_index = self_info[0]["index"]
	}
	log.Println("挑战：" + par_map["index"])
	rival_info, err := db.QueryDuelByIndex(par_map["index"])
	if len(rival_info) == 0 || err != nil {
		log.Println("len(rival_info)", len(rival_info), "err:", err)
		this.ResultString(rw, nil, this.ACTIONERR())
		return
	}
	//目前冒险内容未入，先用自己打自己 注意
	rival_team := this.TeamByUserid(rival_info[0]["userid"])
	//结束
	result_list := this.Duel(self_team, rival_team)
	for k, v := range result_list {
		log.Println(k, v)
	}
	//胜利了
	if strings.Index(result_list[len(result_list)-1], "win") >= 0 {
		err := db.UpdateDuel(par_map["userid"], par_map["username"], par_map["index"])
		if err != nil {
			log.Println(err.Error())
		}
		if self_index != "0" {
			err := db.UpdateDuel(rival_info[0]["userid"], rival_info[0]["username"], self_index)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	this.ResultString(rw, result_list, nil)
}
