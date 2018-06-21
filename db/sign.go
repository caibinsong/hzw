package db

import (
	"../models"
	"fmt"
	"strconv"
	"strings"
)

const (
	DEL_SIGN_SQL  = `delete from sign `
	IS_SIGN_SQL   = `select count(*) as num from usersign where strftime('%%Y-%%m-%%d',sign_time)='%s' and user_id='%s' and sign_type='%s'`
	SIGN_SQL      = `insert into usersign(user_id,sign_type,sign_time) values('%s','%s','%s')`
	SIGN_GIFT_SQL = `select * from sign left join (select count(*) as num from usersign where user_id='%s' and sign_type='%s' and strftime("%%Y-%%m",date('now'))=strftime("%%Y-%%m",sign_time)) a where a.num=sign.sign_count `
)

func IsSign(userid, signtype string) (bool, error) {
	count, err := Count(fmt.Sprintf(IS_SIGN_SQL, models.GetDate(), userid, signtype))
	if err != nil {
		return true, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func Sign(userid, signtype string) error {
	_, err := GetConn().Exec(fmt.Sprintf(SIGN_SQL, userid, signtype, models.GetDate()))
	if err != nil {
		return err
	}
	return nil
}

func SignGift(userid, signtype string) (map[string]string, error) {
	rst, err := Find(fmt.Sprintf(SIGN_GIFT_SQL, userid, signtype))
	if err != nil {
		return rst, err
	}
	_userid, err := strconv.Atoi(userid)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(rst["gift"], "gold:") {
		//金
		_gold, err := strconv.Atoi(subgift(rst["gift"], "gold:"))
		if err != nil {
			return nil, err
		}
		err = AddMoney(_userid, _gold, 0, 0)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(rst["gift"], "diamond:") {
		//钻
		_diamond, err := strconv.Atoi(subgift(rst["gift"], "diamond:"))
		if err != nil {
			return nil, err
		}
		err = AddMoney(_userid, 0, _diamond, 0)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(rst["gift"], "hero:") {
		hero := GetHeroByName(subgift(rst["gift"], "hero:"))
		err = AddHero(userid, hero)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(rst["gift"], "herochip:") {
		hero := GetHeroByName(subgift(rst["gift"], "herochip:"))
		err = AddHero(userid, hero)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(rst["gift"], "eqpt:") {
		eqpt := GetEqptInfoByName(subgift(rst["gift"], "eqpt:"))
		err = AddEqpt(userid, eqpt)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(rst["gift"], "skill:") {
		skill := GetSkillInfoByName(subgift(rst["gift"], "skill:"))
		err = AddSkill(userid, skill)
		if err != nil {
			return nil, err
		}
	}
	return rst, err
}

func subgift(gift, substr string) string {
	return strings.Trim(gift[len(substr):], " ")
}
