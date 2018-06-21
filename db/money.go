package db

import (
	"../models"
	"fmt"
	"log"
	"strconv"
	"sync"
)

var money_Mux sync.Mutex

//在减的时候要先判断是否减后还大于0
func CanUsed(money map[string]string, key string, used int) error {
	num, err := strconv.Atoi(money[key])
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("操作失败")
	}
	if num+used < 0 {
		return fmt.Errorf("no money")
	}
	return nil
}

//通过userid进行消费，大于0为加，小于0为减，0为不操作
func AddMoney(userid, gold, diamond, rmb int) error {
	money_Mux.Lock()
	defer money_Mux.Unlock()
	if gold == 0 && diamond == 0 {
		return fmt.Errorf("无操作")
	}
	user_money := QueryMoneyByUserId(userid)
	if user_money == nil {
		return fmt.Errorf("无此用户")
	}
	sql := "update money set "
	if gold != 0 {
		err := CanUsed(user_money, "gold", gold)
		if err != nil {
			return err
		}
		sql += fmt.Sprintf("gold=gold+%d,", gold)
	}
	if diamond != 0 {
		err := CanUsed(user_money, "diamond", diamond)
		if err != nil {
			return err
		}
		sql += fmt.Sprintf("diamond=diamond+%d,", diamond)
	}
	if rmb > 0 {
		sql += fmt.Sprintf("addrmbsum=addrmbsum+%d,", rmb)
	}
	sql = sql[:len(sql)-1] + fmt.Sprintf(" where user_id=%d", userid)
	_, err := GetConn().Exec(sql)
	if err != nil {
		return err
	}
	AddMoneyLog(userid, gold, diamond, rmb)
	return nil
}

//通过userid查询还有多少
func QueryMoneyByUserId(userid int) map[string]string {
	sql := fmt.Sprintf("select * from money where user_id=%d", userid)
	list, err := Query(sql)
	if err != nil {
		log.Println(sql, err.Error())
		return nil
	}
	if len(list) != 1 {
		log.Println("QueryMoneyByUserId sql:", sql, " not find")
		return nil
	} else {
		return list[0]
	}
}

//添加日志
func AddMoneyLog(userid, gold, diamond, rmb int) {
	sql := fmt.Sprintf("insert into money_log(user_id,add_gold,add_diamond,add_rmb,add_time) values(%d,%d,%d,%d,'%s')",
		userid, gold, diamond, rmb, models.GetTime())
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Println("add money log err:" + err.Error())
	}
	return
}
