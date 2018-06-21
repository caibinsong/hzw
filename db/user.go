package db

import (
	"../models"
	"fmt"
	"log"
	"sync"
)

var user_Mux sync.Mutex

const (
	strength_Default = 50
	strength_Add     = 1
)

//修改登录时间
func UpdateLoginTime(user string) {
	sql := fmt.Sprintf("update user set logintime='%s' where user='%s'", models.GetTime(), user)
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Println("do sql :", sql, " err:", err.Error())
	}
}

//通过用户名查询用户
func FindUserById(userid int) (map[string]string, error) {
	return Find(fmt.Sprintf("select *,strftime('%%Y-%%m-%%d %%H:%%M:%%S',maxstrengthtime) maxstrengthtime from user where id=%d", userid))
}

//用户登录
func UserLogin(user, pwd string) (map[string]string, error) {
	if len(user) < 6 || len(pwd) < 6 {
		return nil, fmt.Errorf("用户名密码不可小于六位")
	}
	if len(user) > 60 || len(pwd) > 60 {
		return nil, fmt.Errorf("用户名密码不可过长")
	}
	sql := fmt.Sprintf("select user.*,gold,diamond,addrmbsum from user left join money on money.user_id=user.id where user='%s' and pwd='%s'", user, models.EnCode(pwd))
	maplist, err := Query(sql)
	if err != nil {
		log.Println("UserLogin error:", err.Error())
		return nil, fmt.Errorf("系统正在维护中")
	}
	if len(maplist) == 0 {
		return nil, fmt.Errorf("用户名密码不正确")
	}
	if len(maplist) != 1 {
		return nil, fmt.Errorf("用户存在问题，请与管理员联系")
	}
	UpdateLoginTime(user)
	return maplist[0], nil
}

//用户注册(先判断用户是否存在，如果不存在则进行注册)
func AddNewUser(user, pwd string) error {
	user_Mux.Lock()
	defer user_Mux.Unlock()
	if len(user) < 6 || len(pwd) < 6 {
		return fmt.Errorf("用户名密码不可小于六位")
	}
	if len(user) > 60 || len(pwd) > 60 {
		return fmt.Errorf("用户名密码不可过长")
	}
	user_map, err := findUser(user)
	if err != nil {
		return fmt.Errorf("用户注册失败，请稍后再试")
	}
	if user_map != nil {
		return fmt.Errorf("用户已存在")
	}
	//添加用户表
	sql := fmt.Sprintf("insert into user(user,pwd,lv,strength,maxstrength,maxstrengthtime,createtime) values('%s','%s',1,%d,%d,'%s','%s')",
		user, models.EnCode(pwd), strength_Default, strength_Default, models.GetTime(), models.GetTime())
	_, err = GetConn().Exec(sql)
	if err != nil {
		log.Println("AddNewUser err:", err.Error(), " sql:", sql)
		return fmt.Errorf("用户注册失败，请稍后再试")
	}
	//查询出id
	usermap, err := findUser(user)
	if err != nil {
		return err
	}
	if usermap == nil {
		return fmt.Errorf("用户注册失败，请稍后再试")
	}
	//添加相关表
	//money表
	sql = fmt.Sprintf("insert into money (user_id,gold,diamond,addrmbsum) values(%s,0,0,0)", usermap["id"])
	_, err = GetConn().Exec(sql)
	if err != nil {
		log.Println("AddNewUser err:", err.Error(), " sql:", sql)
		return fmt.Errorf("用户注册失败，请稍后再试")
	}
	return nil
}

//通过用户名查询用户
func findUser(user string) (map[string]string, error) {
	return Find(fmt.Sprintf("select * from user where user='%s'", user))
}

func IsInitRole(userid, user string) bool {
	sql := fmt.Sprintf("select * from user where id=%s and user='%s'", userid, user)
	result, err := Find(sql)
	if err != nil {
		log.Println(sql, err.Error())
		return false
	}
	if result["id"] != "" && result["name"] == "" {
		return true
	}
	return false
}

//修改名字
func ChangeName(userid, user, name string) error {
	var changename_Mux sync.Mutex
	changename_Mux.Lock()
	defer changename_Mux.Unlock()
	if len(name) > 8 {
		return fmt.Errorf("名字不可过长")
	}
	if name == "" {
		return fmt.Errorf("名字不可为空")
	}
	findname, err := Find(fmt.Sprintf("select * from user where name='%s'", name))
	if err != nil {
		return err
	}
	if findname != nil {
		return fmt.Errorf("名字已存在")
	}
	sql := fmt.Sprintf("update user set name='%s' where id=%s and user='%s'", name, userid, user)
	_, err = GetConn().Exec(sql)
	if err != nil {
		log.Println(err.Error(), sql)
		return fmt.Errorf("请稍后")
	}
	return nil
}

//添加电脑玩家
func AddComputerPlayer(user, pwd, name string) string {
	sql := fmt.Sprintf("insert into user(user,pwd,name,lv,strength,maxstrength,maxstrengthtime,createtime) values('%s','%s','%s',1,%d,%d,'%s','%s')",
		user, models.EnCode(pwd), name, strength_Default, strength_Default, models.GetTime(), models.GetTime())
	_, err := GetConn().Exec(sql)
	if err != nil {
		log.Panic("AddComputerPlayer err:", err.Error(), " sql:", sql)
	}
	usermap, err := findUser(user)
	if err != nil {
		log.Panic(err.Error())
	}
	return usermap["id"]
}
