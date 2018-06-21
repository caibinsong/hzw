package db

import (
	"fmt"
	"log"
)

const (
	duel_hero_count      = "select count(id) as num from duelleaderboard"                            //查询决斗表是否有数据
	duel_delete          = "delete from duelleaderboard"                                             //清空决斗表数据
	add_duel             = "insert into duelleaderboard([index],userid,username) values(%d,%s,'%s')" //添加决斗表数据
	query_duel_by_userid = "select * from duelleaderboard where userid=%s"                           //通过userid查询排名
	query_duel_by_index  = "select * from duelleaderboard where [index]=%s"                          //通过index查询信息
	query_duel_between   = `select duelleaderboard.*, teamheros.heronames
							from duelleaderboard left join (
							select userteam.userid,GROUP_CONCAT(userhero.name) heronames 
							from userteam left join userhero on userteam.userheroid=userhero.id
							group by userteam.userid) teamheros on teamheros.userid=duelleaderboard.userid
							where duelleaderboard.[index]>=%d and duelleaderboard.[index] <=%d`
	update_duel = "update duelleaderboard set userid=%s,username='%s' where [index]=%s"
)

//查询决斗列表有多少个
func DuelHeroCount() int {
	num, err := Count(duel_hero_count)
	if err != nil {
		log.Panic(err.Error())
	}
	return num
}

//删除垃圾数据
func DuelDelete() {
	_, err := GetConn().Exec(duel_delete)
	if err != nil {
		log.Panic(err.Error())
	}
}

//添加决斗表数据
func AddDuel(index int, userid, name string) error {
	_, err := GetConn().Exec(fmt.Sprintf(add_duel, index, userid, name))
	return err
}

//查询自己排名
func QueryDuelByUserid(userid string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(query_duel_by_userid, userid))
}

//通过index查询
func QueryDuelByIndex(index string) ([]map[string]string, error) {
	return Query(fmt.Sprintf(query_duel_by_index, index))
}

//查询自己前面十位玩家的信息
func QueryDuelBetween(start, end int) ([]map[string]string, error) {
	return Query(fmt.Sprintf(query_duel_between, start, end))
}

//修改排名
func UpdateDuel(userid, username, index string) error {
	_, err := GetConn().Exec(fmt.Sprintf(update_duel, userid, username, index))
	return err
}
