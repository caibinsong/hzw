package db

import (
	"fmt"
)

const (
	bag_skill = `select userskill.*,ifnull(userteam.id,9999) sc,userhero.name heroname
from userskill
left join userteam on userteam.skill1=userskill.id or userteam.skill2=userskill.id or userteam.skill3=userskill.id
left join userhero on userhero.id=userteam.userheroid
where userskill.userid=%s
order by sc asc, userskill.quality asc`

	//查询英雄详细信息
	skillinfo = `select * from userskill left join skillinfo on userskill.name=skillinfo.name where userskill.userid=%s and  userskill.id=%s `
)

func QuerySkillInfoById(userid, userskillid string) (map[string]string, error) {
	return Find(fmt.Sprintf(skillinfo, userid, userskillid))
}

//查询英雄基本信息
func QuerySkillBaseInfo() ([]map[string]string, error) {
	maplist, err := Query("select * from skillinfo")
	if err != nil {
		return nil, fmt.Errorf("QuerySkillBaseInfo error:%s", err.Error())
	}
	return maplist, nil
}

func AddSkill(userid string, skill map[string]string) error {
	_, err := GetConn().Exec(fmt.Sprintf("insert into userskill(userid,skillid,name,quality,lv,num) values(%s,%s,'%s','%s',1,%s)",
		userid, skill["id"], skill["name"], skill["quality"], skill["initnum"]))
	return err
}

//背包显示英雄
func BagSkillByUserid(userid string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(bag_skill, userid))
}
