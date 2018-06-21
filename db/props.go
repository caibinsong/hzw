package db

import (
	"../models"
	"fmt"
	"log"
	"sync"
)

const (
	DROP_PROPS_SQL   = `drop table  IF EXISTS [props] `
	QUERY_PROPS_SQL  = `select * from props`
	INSERT_PROPS_SQL = `insert into userprops(userid,propsid,name,quality,lv,num) values('%s','%s','%s','%s','%s','%s')`
	UP_PROPS_SQL     = `update userprops set num=num%s where userid=%s and name='%s' and quality='%s' `
	CANUSEPROPS      = `select count(*) as num from userprops where userid='%s' and name='%s' and quality='%s' and num>=%s`

	//log
	ADD_USE_LOG = `	insert into props_use_log(userid,name,quality,usenum,checkmd5,isused,createtime,usedtime,act)
					values (%s,'%s','%s',%s,'%s',0,'%s','','%s')`

	//log 有些道具用了之后，要再确认这个效果是修改还是放弃
	FIND_USE_LOG      = `select count(*) as num from props_use_log where userid=%s and checkmd5='%s'`
	QUERY_USE_LOG     = `select * from props_use_log where  where userid='%s' and checkmd5='%s' and usedtime=''`
	CONFIRM_PROPS_LOG = `update props_use_log set isused=1, usedtime='%s' where userid='%s' and checkmd5='%s' and usedtime=''`
	UP_USE_LOG        = `update props_use_log set isused=1,usedtime='%s' where userid=%s and checkmd5='%s'`
)

var (
	use_props_mut sync.Mutex
)

func QueryPropsBaseInfo() ([]map[string]string, error) {
	return Query(QUERY_PROPS_SQL)
}

func UseQLZ(userid string, props map[string]string, num, act string) (md5 string, err error) {
	use_props_mut.Lock()
	defer use_props_mut.Unlock()
	if CanUseProps(userid, props["name"], props["quality"], num) {
		md5 = models.GetGuid()
		_, err = GetConn().Exec(fmt.Sprintf(ADD_USE_LOG, userid, props["name"], props["quality"], num, md5, models.GetTime(), act))
		if err != nil {
			log.Println(err)
			return
		}
		sql := fmt.Sprintf(UP_PROPS_SQL, "-"+num, userid, props["name"], props["quality"])
		_, err = GetConn().Exec(sql)
		if err != nil {
			log.Println(sql, err.Error())
		}
	} else {
		err = fmt.Errorf("数量不够")
	}
	return
}

func AddUserProps(userid string, props map[string]string, lv, num string) (err error) {
	sql := ""
	if CanUseProps(userid, props["name"], props["quality"], "0") {
		sql = fmt.Sprintf(UP_PROPS_SQL, "+"+num, userid, props["name"], props["quality"])
	} else {
		sql = fmt.Sprintf(INSERT_PROPS_SQL, userid, props["id"], props["name"], props["quality"], lv, num)
	}
	_, err = GetConn().Exec(sql)
	if err != nil {
		log.Println(sql, err)
	}
	return
}

func CanUseProps(userid, name, quality, num string) bool {
	ok, err := Count(fmt.Sprintf(CANUSEPROPS, userid, name, quality, num))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if ok == 0 {
		return false
	}
	return true
}
