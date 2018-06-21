package db

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

var strength_Mux sync.Mutex
var loc *time.Location = nil

//体力
func queryStrengthBaseById(userid int) (int, int, time.Time, error) {
	if loc == nil {
		loc, _ = time.LoadLocation("Local")
	}
	usre_map, err := FindUserById(userid)
	if err != nil || usre_map == nil {
		return 0, 0, time.Now(), fmt.Errorf("操作失败")
	}
	strength, err := strconv.Atoi(usre_map["strength"])
	if err != nil {
		log.Println("strength", err.Error())
		return 0, 0, time.Now(), fmt.Errorf("请稍后")
	}
	maxstrength, err := strconv.Atoi(usre_map["maxstrength"])
	if err != nil {
		log.Println("maxstrength", err.Error())
		return 0, 0, time.Now(), fmt.Errorf("请稍后")
	}
	maxtime, err := time.ParseInLocation("2006-01-02 15:04:05", usre_map["maxstrengthtime"], loc)
	if err != nil {
		log.Println("maxtime", err.Error())
		return 0, 0, time.Now(), fmt.Errorf("请稍后")
	}
	return strength, maxstrength, maxtime, nil
}

//体力销费，(used<0为加 uesr>0为减)
func UsedStrengthByUserId(userid, used int) error {
	strength_Mux.Lock()
	defer strength_Mux.Unlock()
	//更新当前体力
	hasStrength, err := QueryStrengthByUserId(userid)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("请稍后")
	}
	if hasStrength < used {
		return fmt.Errorf("体力不够")
	}
	//查询最新体力和信息
	strength, maxstrength, maxtime, err := queryStrengthBaseById(userid)
	if err != nil {
		return err
	}

	new_strength := strength - used
	new_time := ""
	if new_strength < maxstrength {
		if strength >= maxstrength {
			//用now()计算最大体力时间
			new_time = time.Now().Add(time.Second * (time.Duration)(60*(maxstrength-new_strength))).Format("2006-01-02 15:04:05")
		} else {
			//给maxtime + 花费的体力时间
			new_time = maxtime.Add(time.Second * (time.Duration)(used*60)).Format("2006-01-02 15:04:05")
		}
	}
	err = updateStrengthByUserId(userid, new_strength, new_time)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

//体力查询
func QueryStrengthByUserId(userid int) (int, error) {
	strength, maxstrength, maxtime, err := queryStrengthBaseById(userid)
	if err != nil {
		return 0, err
	}
	if strength >= maxstrength {
		return strength, nil
	} else {
		now := time.Now()
		if maxtime.Unix() < now.Unix() {
			//体力应该加满了
			log.Println("体力满了")
			if strength < maxstrength {
				strength = maxstrength
			}
		} else {
			log.Println("体力没满")
			//体力没加满，一分钏一个体力计算
			if (maxtime.Unix()-now.Unix())%60 != 0 {
				strength = maxstrength - (int)((maxtime.Unix()-now.Unix())/60+1)
			} else {
				strength = maxstrength - (int)((maxtime.Unix()-now.Unix())/60)
			}
		}
		err = updateStrengthByUserId(userid, strength, "")
		if err != nil {
			return 0, err
		}
		return strength, nil
	}
}

//通过用户ID修改体力
func updateStrengthByUserId(userid, strength int, maxtime string) error {
	sql := ""
	if maxtime == "" {
		sql = fmt.Sprintf("update user set strength=%d where id=%d", strength, userid)
	} else {
		sql = fmt.Sprintf("update user set strength=%d,maxstrengthtime='%s' where id=%d", strength, maxtime, userid)
	}
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Println("err:", sql, err.Error())
		return fmt.Errorf("请稍后")
	}
	return nil
}
