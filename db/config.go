package db

import (
	"../models"
	"fmt"
	"log"
	"sync"
)

var (
	config_mut sync.Mutex
)

const (
	IS_NEWDAY      = `select count(*) num from config where [value]='%s'`
	SET_NEWDAY_SQL = `update config set [value]='%s',[time]='%s' where [key]='newday'`
)

func IsNewDay() bool {
	config_mut.Lock()
	defer config_mut.Unlock()
	sql := fmt.Sprintf(IS_NEWDAY, models.GetDate())
	num, err := Count(sql)
	if err != nil {
		log.Println("do", sql, "err:", err.Error())
	}
	if num == 0 {
		SetNewDay()
		return true
	}
	return false
}

func SetNewDay() {
	sql := fmt.Sprintf(SET_NEWDAY_SQL, models.GetDate(), models.GetTime())
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Println("do", sql, "err:", err.Error())
	}
}
