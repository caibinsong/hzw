package db

import (
	"fmt"
	"log"
)

const (
	//背包显示英雄
	bag_hero = `select userhero.*,heroinfo.quality, ifnull(userteam.id,9999) sc 
from userhero 
left join userteam on userteam.userheroid=userhero.id 
left join heroinfo on userhero.heroid=heroinfo.id
where userhero.userid=%s
order by sc asc, heroinfo.quality asc,lv desc`
	bag_herochip = `select heroinfo.*,userherochip.num from userherochip left join heroinfo on userherochip.heroid=heroinfo.id
where userherochip.userid=%s
order by heroinfo.quality asc,userherochip.num desc`
	//查询英雄详细信息
	heroinfo     = `select * from userhero left join heroinfo on userhero.heroid=heroinfo.id where userhero.id=%s and userhero.userid=%s`
	herochipinfo = `select * from userherochip left join heroinfo on userherochip.heroid=heroinfo.id where userherochip.id=%s and userherochip.userid=%s`
	//用于冒险胜利后计算用的查询
	duelWinHeroInfo = `select userhero.*,
inithp,initatk,initdef,initexp,initwill,hpgrowth,atkgrowth,defgrowth,willgrowth,potgrowth 
from userteam left join userhero on userteam.userheroid=userhero.id
left join heroinfo on heroinfo.id=userhero.heroid
where userteam.[userid]=%s `
	duelWinHero_UpExp = `update userhero set nowexp=nowexp+%s where id=%s`
	duelWinHero_UpLv  = `update userhero set lv=lv+1,hp=hp+%s,atk=atk+%s,def=def+%s,will=will+%s,pot=pot+%s,nowexp=%s where id=%s`
)

func QueryMyHeroInfoById(userid, userheroid string) (map[string]string, error) {
	sql := fmt.Sprintf(heroinfo, userheroid, userid)
	log.Println(sql)
	return Find(sql)
}
func QueryMyHeroChipInfoById(userid, userherochipid string) (map[string]string, error) {
	return Find(fmt.Sprintf(herochipinfo, userherochipid, userid))
}

//查询英雄基本信息
func QueryHeroBaseInfo() ([]map[string]string, error) {
	maplist, err := Query("select * from heroinfo")
	if err != nil {
		return nil, fmt.Errorf("QueryHeroBaseInfo error:%s", err.Error())
	}
	return maplist, nil
}

//背包显示英雄
func BagHeroByUserid(userid string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(bag_hero, userid))
}

//背包显示英雄
func BagHeroChipByUserid(userid string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(bag_herochip, userid))
}

//判断这个用户是否有这个英雄了
func HasHero(userid, heroid string) (bool, error) {
	return Has_Base(fmt.Sprintf("select count(*) as num from userhero where userid=%s and heroid=%s", userid, heroid))
}

//判断这个用户是否有这个英雄碎片了
func HasHeroChip(userid, heroid string) (bool, error) {
	return Has_Base(fmt.Sprintf("select count(*) as num from userherochip where userid=%s and heroid=%s", userid, heroid))
}

func Has_Base(sql string) (bool, error) {
	maplist, err := Query(sql)
	if err != nil {
		return false, err
	}
	if maplist[0]["num"] == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func AddHero(userid string, hero map[string]string) error {
	hashero, err := HasHero(userid, hero["id"])
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if hashero {
		haschip, err := HasHeroChip(userid, hero["id"])
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if haschip {
			err = updateHerochip(userid, hero, 5)
		} else {
			err = addHerochip(userid, hero, 5)
		}
	} else {
		err = addHero(userid, hero)
	}
	return err
}

//添加英雄碎片
func addHerochip(userid string, hero map[string]string, num int) error {
	_, err := GetConn().Exec(fmt.Sprintf("insert into userherochip(userid,heroid,name,num) values(%s,%s,'%s',%d)", userid, hero["id"], hero["name"], num))
	return err
}

//update英雄碎片
func updateHerochip(userid string, hero map[string]string, num int) error {
	_, err := GetConn().Exec(fmt.Sprintf("update userherochip set num=num+%d where userid=%s and heroid=%s and name='%s'", num, userid, hero["id"], hero["name"]))
	return err
}

//在没有英雄的前提下添加英雄
func addHero(userid string, hero map[string]string) error {
	sql := fmt.Sprintf(`insert into  userhero(userid,heroid,name,lv,hp,atk,def,will,pot,nowexp,newselfskill,selfskilllv) values(%s,%s,'%s',1,%s,%s,%s,%s,%s,0,'%s',1)`,
		userid, hero["id"], hero["name"], hero["inithp"], hero["initatk"], hero["initdef"], hero["initwill"], hero["potgrowth"], hero["selfskill"])
	_, err := GetConn().Exec(sql)
	return err
}

func DuelWinHeroInfo(userid string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(duelWinHeroInfo, userid))
}

func DuelWinHero_UpExp(add_exp, userheroid string) error {
	sql := fmt.Sprintf(duelWinHero_UpExp, add_exp, userheroid)
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Println(err.Error(), sql)
	}
	return err
}

func DuelWinHero_UpLv(add_hp, add_atk, add_def, add_will, add_pot, nowexp, userheroid string) error {
	sql := fmt.Sprintf(duelWinHero_UpLv, add_hp, add_atk, add_def, add_will, add_pot, nowexp, userheroid)
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Println(err.Error(), sql)
	}
	return err
}
