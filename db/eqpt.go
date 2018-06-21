package db

import (
	"fmt"
)

const (
	bag_eqpt = `select usereqpt.*,ifnull(userteam.id,9999) sc,eqptinfo.typename,userhero.name heroname from usereqpt
left join eqptinfo on eqptinfo.id=usereqpt.eqptid
left join userteam on userteam.%s=usereqpt.id 
left join userhero on userteam.userheroid=userhero.id
where trim(eqptinfo.typename,' ')='%s' and usereqpt.userid=%s
order by sc asc, usereqpt.quality asc`
	//查询装备详细信息
	eqptinfo = `select * from usereqpt left join eqptinfo on usereqpt.eqptid=eqptinfo.id where usereqpt.userid=%s and usereqpt.id=%s`
)

func QueryEqptInfoById(userid, usereqptid string) (map[string]string, error) {
	return Find(fmt.Sprintf(eqptinfo, userid, usereqptid))
}

//查询英雄基本信息
func QueryEqptBaseInfo() ([]map[string]string, error) {
	maplist, err := Query("select * from eqptinfo")
	if err != nil {
		return nil, fmt.Errorf("QueryEqptBaseInfo error:%s", err.Error())
	}
	return maplist, nil
}

func AddEqpt(userid string, eqpt map[string]string) error {
	_, err := GetConn().Exec(fmt.Sprintf("insert into usereqpt(userid,eqptid,name,quality,lv,num) values (%s,%s,'%s','%s',1,%s)",
		userid, eqpt["id"], eqpt["name"], eqpt["quality"], eqpt["initnum"]))
	return err
}

//背包显示英雄
func BagEqptByUserid(userid, types string) ([]map[string]string, error) {
	zd := "eqpt"
	if types == "武器" {
		zd += "1"
	}
	if types == "防具" {
		zd += "2"
	}
	if types == "饰品" {
		zd += "3"
	}
	return Query(fmt.Sprintf(bag_eqpt, zd, types, userid))
}
