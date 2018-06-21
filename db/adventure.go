package db

import (
	"fmt"
	"log"
)

const (
	_QUERY_ADVENTURE = `select adventure.*,ifnull(useradventure.id,0) has,ifnull(useradventure.todaytimes,0) todaytimes
from adventure left join useradventure on useradventure.[adventureindex]=adventure.[index] and useradventure.userid=%s
order by adventure.[index] asc`
	_QUERY_ADVENTURE_INFO    = `select * from adventure where [index]=%s`
	_QUERY_ADVENTURE_HERO    = `select * from adventurehero where adventureindex=%s`
	_NEED_ADD_ADVENTURE_INFO = `select count(*) num from useradventure where userid=%s and adventureindex=%s`
	_UPDATE_ADVENTURE        = `update useradventure set todaytimes=todaytimes+1 where userid=%s and adventureindex=%s`
	_INSERT_ADVENTURE        = `insert into useradventure(userid,adventureindex,todaytimes) values(%s,%s,1)`
	_NEWDAY                  = `update useradventure set todaytimes=0`
)

//查询关卡类型 在最后一关加个非0（因为没打过的关 下一关是要可以打的）
func QueryAdventure(userid string) ([]map[string]string, error) {
	rst, err := Query(fmt.Sprintf(_QUERY_ADVENTURE, userid))
	if err != nil {
		return rst, err
	}
	for k, _ := range rst {
		if rst[k]["has"] == "0" {
			rst[k]["has"] = "-1"
			break
		}
	}
	return rst, nil
}

func QueryAdventureInfo(index string) (map[string]string, error) {
	return Find(fmt.Sprintf(_QUERY_ADVENTURE_INFO, index))
}

//查询关卡英雄
func QueryAdventureHero(adventureindex string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(_QUERY_ADVENTURE_HERO, adventureindex))
}

func UpdateAdventure(userid, adventureindex string) error {
	count, err := Count(fmt.Sprintf(_NEED_ADD_ADVENTURE_INFO, userid, adventureindex))
	if err != nil {
		return err
	}
	sql := ""
	if count == 1 {
		sql = fmt.Sprintf(_UPDATE_ADVENTURE, userid, adventureindex)
	} else {
		sql = fmt.Sprintf(_INSERT_ADVENTURE, userid, adventureindex)
	}
	_, err = GetConn().Exec(sql)
	return err
}

func AdventureNewDay() {
	_, err := GetConn().Exec(_NEWDAY)
	if err != nil {
		log.Println(err.Error())
	}
}
