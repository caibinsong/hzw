package controllers

import (
	"../db"
	"../models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BaseConttroller struct {
}

var coder = base64.StdEncoding
var xor []byte = []byte{14, 3, 6, 15}
var pMap map[string]string = map[string]string{"血量": "HP", "攻击": "ATK", "防御": "DEF", "意志": "WILL"}

const (
	BAOJI     = 1.2
	BAOJI_STR = "1.2"
	SKILL     = 1.2
)

//接收
func (this *BaseConttroller) GetParameter(req *http.Request) map[string]string {
	rst := make(map[string]string)
	////后面要删除掉
	req.ParseForm()
	for k, v := range req.Form {
		rst[k] = v[0]
	}
	//
	body_b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
		return rst
	}
	arr := strings.Split(this.Http_DeCode(string(body_b)), "<>")
	index := -1
	for _, one := range arr {
		index = strings.Index(one, "=")
		if index > 0 {
			//log.Println(one[:index], one[index+1:])
			rst[one[:index]] = one[index+1:]
		}
	}
	return rst
}

//返回 最后加时间为了生成加密后，每次返回的信息都不一样
func (this *BaseConttroller) ResultString(rw http.ResponseWriter, msg interface{}, msg_err error) {
	if msg_err == nil {
		rst_msg := ""
		if msg != nil {
			jsondata, err := json.Marshal(msg)
			if err != nil {
				log.Println(err.Error())
				io.WriteString(rw, "err:操作失败，请稍后 END:"+models.GetTime())
				return
			}
			rst_msg = string(jsondata)
			if rst_msg != "" {
				rst_msg = ":" + rst_msg
			}
		}
		io.WriteString(rw, this.Http_EnCode("success"+fmt.Sprint(rst_msg)+" END:"+models.GetTime()))
	} else {
		io.WriteString(rw, this.Http_EnCode("err:"+msg_err.Error()+" END:"+models.GetTime()))
	}
}

//http通信时参数内容的解密
func (this *BaseConttroller) Http_DeCode(str string) string {
	return Decode(str, 1)
}

//http通信返回时的加密方法
func (this *BaseConttroller) Http_EnCode(str string) string {
	return Encoder(str, 1)
}

func (this *BaseConttroller) ACTIONERR() error {
	return fmt.Errorf("操作失败")
}

func Encoder(str string, num int) string {
	if num < 1 || num > 4 {
		return str
	}
	bts := []byte(str)
	xor_num := 0
	for i := 0; i < len(bts); i++ {
		bts[i] = bts[i] ^ xor[xor_num]
		xor_num++
		if xor_num >= num {
			xor_num = 0
		}
	}
	return coder.EncodeToString(bts)
}

func Decode(str string, num int) string {
	if num < 1 || num > 4 {
		return str
	}
	xor_num := 0
	b, _ := coder.DecodeString(str)
	for i := 0; i < len(b); i++ {
		b[i] = b[i] ^ xor[xor_num]
		xor_num++
		if xor_num >= num {
			xor_num = 0
		}
	}
	return string(b)
}

func (this *BaseConttroller) AddIntNum(s_num0, s_num1 string) int {
	num0, err := strconv.Atoi(s_num0)
	if err != nil {
		num0 = 0
	}
	num1, err := strconv.Atoi(s_num1)
	if err != nil {
		num1 = 0
	}
	num0 += num1
	return num0
}

//加减
func (this *BaseConttroller) Add(s_num0, s_num1 string) string {
	var num0, num1 float64 = 0xFFFFFFF, 0xFFFFFFF
	var err error
	num0, err = strconv.ParseFloat(s_num0, 64)
	if err != nil {
		num0 = 0
	}
	num1, err = strconv.ParseFloat(s_num1, 64)
	if err != nil {
		num1 = 0
	}
	return strconv.FormatFloat(num0+num1, 'f', 4, 64)
}

//str 转float64
func (this *BaseConttroller) str2float(s_num0 string) float64 {
	var num0 float64 = 0x000000
	var err error
	num0, err = strconv.ParseFloat(s_num0, 64)
	if err != nil {
		num0 = 0
	}
	return num0
}

func (this *BaseConttroller) ForMatPercent(s_num0 string) string {
	var num0 float64 = 0xFFFFFFF
	var err error
	num0, err = strconv.ParseFloat(s_num0, 64)
	if err != nil {
		num0 = 0
	}
	num0 *= 100
	return strconv.FormatFloat(num0, 'f', 4, 64)
}

//乘
func (this *BaseConttroller) MultiplyIntNum(s_num0 string, s_num1 string) int {
	num0, err := strconv.Atoi(s_num0)
	if err != nil {
		num0 = 1
	}
	num1, err := strconv.Atoi(s_num1)
	if err != nil {
		num1 = 1
	}
	num0 = num0 * num1
	return num0
}

//乘
func (this *BaseConttroller) Multiply(s_num string, f_num float64) string {
	return strconv.FormatFloat(this.str2float(s_num)*f_num, 'f', 2, 64)
}

//提取缘名 :前面的为缘名 否则返回全部
func (this *BaseConttroller) GetLuckName(str string) string {
	index := strings.Index(str, ":")
	if index > 0 {
		return strings.TrimSpace(str[:index])
	} else {
		return str
	}
}

//提取激活缘所需英雄名
func (this *BaseConttroller) GetLuckHeroName(str string) string {
	str = strings.Replace(strings.TrimSpace(str), " :", ":", 1)
	start := strings.Index(str, " ")
	end := strings.LastIndex(str, " ")
	if start == end || start < 0 || end < 0 {
		return str
	}
	return str[start+1 : end]
}

//提取缘添加属性的值
func (this *BaseConttroller) GetAddNum(str string) (string, float64) {
	index := strings.LastIndex(str, ",")
	if index < 0 {
		return "", 0
	}
	str = str[index+1:]
	arr := strings.Split(str, "提升")
	if len(arr) == 2 {
		var num float64 = 0xFFFFFFF
		var err error
		arr[1] = strings.TrimSpace(strings.Replace(arr[1], "%", "", 1))
		num, err = strconv.ParseFloat(arr[1], 64)
		if err != nil {
			log.Println(err.Error())
			return pMap[arr[0]], 0
		}
		return pMap[arr[0]], (float64)(num / (float64)(100.0))
	} else {
		return "", 0
	}
}

//判断缘份是否激活
func (this *BaseConttroller) ActivateLuck(elements string, team []map[string]string, team_k int) bool {
	if strings.TrimSpace(elements) == "" {
		return false
	}
	arrElement := strings.Split(elements, "、")
	for eleindex, element := range arrElement {
		for k, hero := range team {
			if k == team_k {
				if hero["newselfskill"] == element || hero["skill1_name"] == element ||
					hero["skill2_name"] == element || hero["skill3_name"] == element ||
					hero["eqpt1_name"] == element || hero["eqpt1_name"] == element || hero["eqpt1_name"] == element {
					arrElement[eleindex] = "ok"
					break
				}
			} else {
				if hero["name"] == element {
					arrElement[eleindex] = "ok"
					break
				}
			}
		}
	}
	for _, isok := range arrElement {
		if isok != "ok" {
			return false
		}
	}
	return true
}

//通过userid查询阵容，并计算出装备技能缘份加成后的属性
func (this *BaseConttroller) TeamByUserid(userid string) []map[string]string {
	rst, err := db.QueryMyTeam(userid)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	//计算攻防和缘
	for k, _ := range rst {
		//百分比加成(基数为1开始，可直接乘)
		add_p := make(map[string]float64)
		add_p["HP"] = 1.0
		add_p["ATK"] = 1.0
		add_p["DEF"] = 1.0
		add_p["WILL"] = 1.0
		//把装备加成加上
		rst[k]["atk"] = this.Add(rst[k]["atk"], rst[k]["eqpt1_num"])
		rst[k]["def"] = this.Add(rst[k]["def"], rst[k]["eqpt2_num"])
		rst[k]["hp"] = this.Add(rst[k]["hp"], rst[k]["eqpt3_num"])
		//把技能加成加到百分比加成上
		if rst[k]["skill0typename"] == "ATK" || rst[k]["skill0typename"] == "DEF" || rst[k]["skill0typename"] == "WILL" ||
			rst[k]["skill0typename"] == "HP" {
			add_p[rst[k]["skill0typename"]] = add_p[rst[k]["skill0typename"]] + this.str2float(rst[k]["skill0_num"]) +
				this.str2float(rst[k]["skill0growthnum"])*(this.str2float(rst[k]["selfskilllv"])-1)
		}
		if rst[k]["skill1typename"] == "ATK" || rst[k]["skill1typename"] == "DEF" || rst[k]["skill1typename"] == "WILL" ||
			rst[k]["skill1typename"] == "HP" {
			add_p[rst[k]["skill1typename"]] += this.str2float(rst[k]["skill1_num"])
		}
		if rst[k]["skill2typename"] == "ATK" || rst[k]["skill2typename"] == "DEF" || rst[k]["skill2typename"] == "WILL" ||
			rst[k]["skill2typename"] == "HP" {
			add_p[rst[k]["skill2typename"]] += this.str2float(rst[k]["skill2_num"])
		}
		if rst[k]["skill3typename"] == "ATK" || rst[k]["skill3typename"] == "DEF" || rst[k]["skill3typename"] == "WILL" ||
			rst[k]["skill3typename"] == "HP" {
			add_p[rst[k]["skill3typename"]] += this.str2float(rst[k]["skill3_num"])
		}
		//提取缘信息
		for i := 1; i <= 6; i++ {
			rst[k][fmt.Sprintf("luck%dname", i)] = this.GetLuckName(rst[k][fmt.Sprintf("luck%d", i)])
			rst[k][fmt.Sprintf("luck%dheroname", i)] = this.GetLuckHeroName(rst[k][fmt.Sprintf("luck%d", i)])
			ok := this.ActivateLuck(rst[k][fmt.Sprintf("luck%dheroname", i)], rst, k)
			rst[k][fmt.Sprintf("luck%djh", i)] = fmt.Sprint(ok)
			if ok {
				addtype, addnum := this.GetAddNum(rst[k][fmt.Sprintf("luck%d", i)])
				add_p[addtype] += float64(addnum)
			}
		}
		//计算最后的属性
		rst[k]["hp"] = this.Multiply(rst[k]["hp"], add_p["HP"])
		rst[k]["atk"] = this.Multiply(rst[k]["atk"], add_p["ATK"])
		rst[k]["def"] = this.Multiply(rst[k]["def"], add_p["DEF"])
		rst[k]["will"] = this.Multiply(rst[k]["will"], add_p["WILL"])
	}
	return rst
}

//决斗 传入两个阵容
func (this *BaseConttroller) Duel(self_team, rival_team []map[string]string) []string {
	rand.Seed(time.Now().UnixNano() / 100)
	result_list := make([]string, 0)
	self_team_str := "self_team_hero:"
	rival_team_str := "rival_team_hero:"
	for _, hero := range self_team {
		self_team_str = self_team_str + hero["name"] + " "
	}
	for _, hero := range rival_team {
		rival_team_str = rival_team_str + hero["name"] + " "
	}
	result_list = append(result_list, self_team_str)
	result_list = append(result_list, rival_team_str)
	self_team_sum_hp, self_team_sum_will, rival_team_sum_hp, rival_team_sum_will, isend, endmsg := this.SUM_HP_WILL2(self_team, rival_team)
	log.Println(self_team_sum_will, rival_team_sum_will)

	//判断是否分出胜负
	if isend {
		result_list = append(result_list, endmsg)
	}

	//判断谁先攻击
	selfatk := self_team_sum_will >= rival_team_sum_will
	self_team[0]["teamname"] = "self_team"
	rival_team[0]["teamname"] = "rival_team"

	//开始
	for huihe := 1; huihe <= 2; huihe++ {
		result_list = append(result_list, fmt.Sprintf("第%d回合", huihe))
		if selfatk {
			self_team, rival_team, result_list = this.attack(0, self_team, rival_team, result_list)
			rival_team, self_team, result_list = this.attack(0, rival_team, self_team, result_list)
			self_team, rival_team, result_list = this.attack(0, self_team, rival_team, result_list)
			rival_team, self_team, result_list = this.attack(0, rival_team, self_team, result_list)
		} else {
			rival_team, self_team, result_list = this.attack(0, rival_team, self_team, result_list)
			self_team, rival_team, result_list = this.attack(0, self_team, rival_team, result_list)
			rival_team, self_team, result_list = this.attack(0, rival_team, self_team, result_list)
			self_team, rival_team, result_list = this.attack(0, self_team, rival_team, result_list)
		}
		log.Println(huihe, "q------HP0Remove--------self", len(self_team))
		self_team = this.HP0Remove(self_team)
		log.Println(huihe, "q------HP0Remove--------self", len(rival_team))
		rival_team = this.HP0Remove(rival_team)
		selfatk = !selfatk
		if len(self_team) == 0 || len(rival_team) == 0 {
			break
		}
		log.Println(huihe, "h------HP0Remove--------self", len(self_team))
		for k, hero := range self_team {
			log.Println(fmt.Sprintf("%d:%s %s hp:%s atk:%s def:%s will:%s", k, hero["name"], hero["teamindex"], hero["hp"], hero["atk"], hero["def"], hero["will"]))
		}

		log.Println(huihe, "h------HP0Remove--------rival", len(rival_team))
		for k, hero := range rival_team {
			log.Println(fmt.Sprintf("%d:%s %s hp:%s atk:%s def:%s will:%s", k, hero["name"], hero["teamindex"], hero["hp"], hero["atk"], hero["def"], hero["will"]))
		}
	}
	if len(self_team) != 0 && len(rival_team) != 0 {
		result_list = append(result_list, "第3回合")
	}
	self_team_sum_hp, self_team_sum_will = this.SUM_HP_WILL(self_team)
	rival_team_sum_hp, rival_team_sum_will = this.SUM_HP_WILL(rival_team)
	if (self_team_sum_hp + self_team_sum_will) >= (rival_team_sum_hp + rival_team_sum_will) {
		result_list = append(result_list, "self win")
	} else {
		result_list = append(result_list, "self fail")
	}
	return result_list
}

//随机数 0 - maxnum之间的
func (this *BaseConttroller) Rand(maxnum int) int {
	return rand.Intn(maxnum)
}

func (this *BaseConttroller) skillTypeNameIsADorAOE(typename string) bool {
	if typename == "AD" || typename == "AOE" {
		return true
	} else {
		return false
	}
}

//计算攻击类型 0-30 普攻 ;30-40普攻暴击; 40-55 0 ; 55-70 1 ; 70-85 2 ; 85-100 3
func (this *BaseConttroller) GetAtkType(hero map[string]string) (int, string, bool, float64) {
	randnum := this.Rand(100)
	isAOE := false
	if randnum < 30 {
		return -2, "", isAOE, 1
	}
	if randnum < 40 {
		return -1, "暴击", isAOE, BAOJI
	}
	skill_index := 0
	skillname := ""
	var mutil float64 = 0x0000000
	mutil = 1

	if randnum < 55 {
		if this.skillTypeNameIsADorAOE(hero["skill0typename"]) {
			//用技能攻击
			skill_index = 0
			skillname = hero["newselfskill"]
			if hero["skill0typename"] == "AOE" {
				isAOE = true
			}
			growthnum := this.Multiply(this.Add(hero["selfskilllv"], "-1"), this.str2float(hero["skill0growthnum"]))
			mutil = this.str2float(this.Add(hero["skill0_num"], growthnum))/100 + BAOJI
		} else {
			//无技能或是加成类技能，变成普攻
			return -2, "", isAOE, 1
		}
	} else if randnum < 70 {
		if this.skillTypeNameIsADorAOE(hero["skill1typename"]) {
			skill_index = 1
			skillname = hero["skill1_name"]
			if hero["skill1typename"] == "AOE" {
				isAOE = true
			}
			mutil = this.str2float(hero["skill1_num"])/100 + BAOJI
		} else {
			//无技能或是加成类技能，变成普攻
			return -2, "", isAOE, 1
		}
	} else if randnum < 85 {
		if this.skillTypeNameIsADorAOE(hero["skill2typename"]) {
			skill_index = 2
			skillname = hero["skill2_name"]
			if hero["skill2typename"] == "AOE" {
				isAOE = true
			}
			mutil = this.str2float(hero["skill2_num"])/100 + BAOJI
		} else {
			//无技能或是加成类技能，变成普攻
			return -2, "", isAOE, 1
		}
	} else {
		if this.skillTypeNameIsADorAOE(hero["skill3typename"]) {
			skill_index = 3
			skillname = hero["skill3_name"]
			if hero["skill3typename"] == "AOE" {
				isAOE = true
			}
			mutil = this.str2float(hero["skill3_num"])/100 + BAOJI
		} else {
			//无技能或是加成类技能，变成普攻
			return -2, "", isAOE, 1
		}
	}
	return skill_index, skillname, isAOE, mutil
}

func (this *BaseConttroller) SUM_HP_WILL2(self_team, rival_team []map[string]string) (float64, float64, float64, float64, bool, string) {
	self_hp, self_will := this.SUM_HP_WILL(self_team)
	rival_hp, rival_will := this.SUM_HP_WILL(rival_team)
	isEnd := false
	endMsg := ""
	if self_hp == 0 {
		isEnd = true
		endMsg = "self fail"
	}
	if rival_hp == 0 {
		isEnd = true
		endMsg = "self win"
	}
	return self_hp, self_will, rival_hp, rival_will, isEnd, endMsg
}

//计算出sumHP和sumWILL 用于判断对面是否全部阵亡和谁先攻击
func (this *BaseConttroller) SUM_HP_WILL(team []map[string]string) (float64, float64) {
	var sum_hp, sum_will float64 = 0x0000000, 0x0000000
	for _, hero := range team {
		//log.Println(fmt.Sprintf("%d:%s %s hp:%s atk:%s def:%s will:%s", k, hero["name"], hero["teamindex"], hero["hp"], hero["atk"], hero["def"], hero["will"]))
		sum_hp += this.str2float(hero["hp"])
		sum_will += this.str2float(hero["will"])
	}
	return sum_hp, sum_will
}

//攻击
func (this *BaseConttroller) attack(attack_index int, atk_team, def_team []map[string]string, result_list []string) ([]map[string]string, []map[string]string, []string) {
	if attack_index >= len(atk_team) && attack_index >= len(def_team) {
		return atk_team, def_team, result_list
	}
	if attack_index >= len(atk_team) {
		return this.attack(attack_index+1, def_team, atk_team, result_list)
	}
	if this.str2float(atk_team[attack_index]["hp"]) != 0 {
		def_index := this.FindDefObject(attack_index, def_team)
		if def_index != -1 {
			//算攻击类型 是普攻还是技能 技能分AOE和非AOE
			/*atktype*/ _, skillname, isAOE, multiple := this.GetAtkType(atk_team[attack_index])

			if isAOE {
				//aoe技能
				skillname = "使用AOE技能 " + skillname
				bei_gong_ji_Str := ""
				for index, def_hero := range def_team {
					if this.str2float(def_hero["hp"]) == 0 {
						continue
					}
					hurt := this.hurt(atk_team[attack_index], def_hero, multiple)
					def_team[index]["hp"] = this.Add(def_team[index]["hp"], hurt)
					bei_gong_ji_Str += fmt.Sprintf("%d %s 伤害为 %s 剩 %s;", index, def_team[index]["name"], hurt, def_team[index]["hp"])

				}
				result_list = append(result_list,
					fmt.Sprintf("%s %d %s %s 攻击 (%s)",
						atk_team[0]["teamname"], attack_index, atk_team[attack_index]["name"], skillname,
						bei_gong_ji_Str))
			} else {
				if skillname == "" {
					//普攻
					skillname = "普攻"
				} else if skillname == "暴击" {
					//普攻暴击
					skillname = "普攻暴击"
				} else {
					//单攻技能
					skillname = "使用技能 " + skillname
				}
				hurt := this.hurt(atk_team[attack_index], def_team[def_index], multiple)
				def_team[def_index]["hp"] = this.Add(def_team[def_index]["hp"], hurt)
				result_list = append(result_list,
					fmt.Sprintf("%s %d %s %s 攻击 %d %s 伤害为 %s 剩 %s", atk_team[0]["teamname"], attack_index, atk_team[attack_index]["name"], skillname,
						def_index, def_team[def_index]["name"], hurt, def_team[def_index]["hp"]))
			}
			for index, def_hero := range def_team {
				if this.str2float(def_hero["hp"]) == 0 && def_hero["isdie"] != "1" {
					def_team[index]["isdie"] = "1"
					result_list = append(result_list, fmt.Sprintf("%s %d %s 阵亡", def_team[0]["teamname"], index, def_team[index]["name"]))
				}
			}
			//log.Println(atk_team[0]["teamname"], attack_index, atk_team[attack_index]["name"], "攻击", def_index,
			//	def_team[def_index]["name"], "伤害为", hurt, "剩", def_team[def_index]["hp"])
		}
	}
	//
	return this.attack(attack_index+1, def_team, atk_team, result_list)
}

//查找攻击对象
func (this *BaseConttroller) FindDefObject(atk_index int, def_team []map[string]string) int {
	if atk_index > len(def_team) {
		atk_index = len(def_team)
	}
	def_index := -1
	for i := atk_index; i < len(def_team); i++ {
		if this.str2float(def_team[i]["hp"]) != 0 {
			def_index = i
			break
		}
	}
	if def_index == -1 {
		for i := 0; i < atk_index; i++ {
			if this.str2float(def_team[i]["hp"]) != 0 {
				def_index = i
				break
			}
		}
	}
	return def_index
}

//计算伤害
func (this *BaseConttroller) hurt(atk_team, def_team map[string]string, multiple float64) string {
	atk := this.str2float(atk_team["atk"])
	def := this.str2float(def_team["def"])
	hp := this.str2float(def_team["hp"])
	_hurt := atk*multiple - def
	if _hurt <= 0 {
		_hurt = 1
	}
	if _hurt >= hp {
		_hurt = hp
	}
	_hurt = _hurt * -1
	return strconv.FormatFloat(_hurt, 'f', 4, 64)
}

//把HP为0的移掉
func (this *BaseConttroller) HP0Remove(team []map[string]string) []map[string]string {
	log.Println("HP0Remove:", len(team))
	for i := 0; i < len(team); {
		log.Println(len(team))
		hp := this.str2float(team[i]["hp"])
		if hp == 0 {
			log.Println(i, team[i]["name"], team[i]["hp"], hp, "remove")
			team = append(team[:i], team[i+1:]...)
		} else {
			log.Println(i, team[i]["name"], team[i]["hp"], hp, "no remove")
			i++
		}
	}
	return team
}
