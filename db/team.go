package db

import (
	"fmt"
	"log"
)

const (
	myteam_SQL = `select userteam.id,userteam.teamindex,userteam.userheroid
,userhero.heroid,userhero.name,userhero.lv,userhero.hp,userhero.atk,userhero.def,userhero.will
,userhero.newselfskill,userhero.selfskilllv,skillinfo0.typename as skill0typename,skillinfo0.initnum as skill0_num,skillinfo0.growthnum as skill0growthnum
,heroinfo.luck1,heroinfo.luck2,heroinfo.luck3,heroinfo.luck4,heroinfo.luck5,heroinfo.luck6
,skill1.id as skill1_skillid,skill1.name as skill1_name,skill1.num as skill1_num,skillinfo1.typename as skill1typename
,skill2.id as skill2_skillid,skill2.name as skill2_name,skill2.num as skill2_num,skillinfo2.typename as skill2typename
,skill3.id as skill3_skillid,skill3.name as skill3_name,skill3.num as skill3_num,skillinfo3.typename as skill3typename
,eqpt1.id as eqpt1_eqptid,eqpt1.name as eqpt1_name,eqpt1.num as eqpt1_num
,eqpt2.id as eqpt2_eqptid,eqpt2.name as eqpt2_name,eqpt2.num as eqpt2_num
,eqpt3.id as eqpt3_eqptid,eqpt3.name as eqpt3_name,eqpt3.num as eqpt3_num  
from userteam 
left join userhero on userhero.id=userteam.userheroid
left join heroinfo on heroinfo.id=userhero.heroid
left join userskill skill1 on skill1.id=userteam.skill1
left join userskill skill2 on skill2.id=userteam.skill2
left join userskill skill3 on skill3.id=userteam.skill3
left join skillinfo skillinfo0 on skillinfo0.name=userhero.newselfskill
left join skillinfo skillinfo1 on skillinfo1.id=skill1.skillid
left join skillinfo skillinfo2 on skillinfo2.id=skill2.skillid
left join skillinfo skillinfo3 on skillinfo3.id=skill3.skillid
left join usereqpt eqpt1 on eqpt1.id=userteam.eqpt1
left join usereqpt eqpt2 on eqpt2.id=userteam.eqpt2
left join usereqpt eqpt3 on eqpt3.id=userteam.eqpt3
where userteam.userid=%s order by userteam.teamindex`
)

//查询阵容
func QueryMyTeam(userid string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(myteam_SQL, userid))
}

//设置阵容
func SetHero(userid, userheroid, teamindex string) error {
	//查询当前阵容
	team, err := QueryMyTeam(userid)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("请稍后")
	}
	//判断是否已经上阵
	flag := -1
	for k, one_hero := range team {
		if one_hero["userheroid"] == userheroid {
			return fmt.Errorf("已上阵")
		}
		if one_hero["teamindex"] == teamindex {
			flag = k
		}
	}
	sql := fmt.Sprintf("insert into userteam(userid,userheroid,teamindex,skill1,skill2,skill3,eqpt1,eqpt2,eqpt3) values(%s,%s,%s,0,0,0,0,0,0)",
		userid, userheroid, teamindex)
	//上阵
	if flag != -1 {
		sql = fmt.Sprintf("update userteam set userheroid=%s where id=%s", userheroid, team[flag]["id"])
	}
	_, err = GetConn().Exec(sql)
	if err != nil {
		log.Println(team, flag, sql)
		log.Println(err.Error())
		return fmt.Errorf("请稍后")
	}
	return nil
}

func SetEqptSkill(userid, teamindex, types, usereqptid string) error {
	//查询当前阵容
	team, err := QueryMyTeam(userid)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("请稍后")
	}
	//判断是否已经上阵
	flag := -1
	_type := types[:len(types)-1]
	for k, one_hero := range team {
		if (one_hero[_type+"1_"+_type+"id"] == usereqptid || one_hero[_type+"2_"+_type+"id"] == usereqptid || one_hero[_type+"3_"+_type+"id"] == usereqptid) && usereqptid != "0" {
			return fmt.Errorf("已使用")
		}
		if one_hero["teamindex"] == teamindex {
			flag = k
		}
	}
	//上阵
	if flag == -1 {
		return fmt.Errorf("人物未上场")
	}
	sql := fmt.Sprintf("update userteam set %s=%s where id=%s", types, usereqptid, team[flag]["id"])
	_, err = GetConn().Exec(sql)
	if err != nil {
		log.Println(sql, err.Error())
		return fmt.Errorf("请稍后")
	}
	return nil
}
