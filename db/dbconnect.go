package db

import (
	//"../models"
	"database/sql"
	"fmt"
	_ "gosql/sqlite3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	//"strings"
)

var database string = "hzw.db" //数据库名
var dbconn *sql.DB = nil
var Hero_Base_Info []map[string]string  //基础英雄信息
var Eqpt_Base_Info []map[string]string  //基础装备信息
var Skill_Base_Info []map[string]string //基础技能信息
var Props_Base_Info []map[string]string //基础道具信息

func GetConn() *sql.DB {
	return dbconn
}

//打开数据库
func DBConn() error {
	sqldb, err := sql.Open("sqlite3", database)
	if err != nil {
		return fmt.Errorf("打开数据库失败：%s", err.Error())
	}
	dbconn = sqldb
	err = CreateTable()
	if err != nil {
		return err
	}
	return Cache()
}

//建表
func CreateTable() error {
	err := filepath.Walk("./table", func(path string, fs os.FileInfo, err error) error {
		if !fs.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("create table ,read sql file: %s ,error:%s", path, err.Error())
			}
			_, err = dbconn.Exec(string(b))
			if err != nil {
				return fmt.Errorf("create table , file: %s,exec:%s ,error:%s", path, string(b), err.Error())
			}
		}
		return nil
	})
	return err
}

//cache
func Cache() error {
	var err error = nil
	Hero_Base_Info, err = QueryHeroBaseInfo()
	if err != nil {
		return err
	}
	if len(Hero_Base_Info) == 0 {
		return fmt.Errorf("英雄库不可为空")
	}
	Eqpt_Base_Info, err = QueryEqptBaseInfo()
	if err != nil {
		return err
	}
	if len(Eqpt_Base_Info) == 0 {
		return fmt.Errorf("装备库不可为空")
	}
	Skill_Base_Info, err = QuerySkillBaseInfo()
	if err != nil {
		return err
	}
	if len(Skill_Base_Info) == 0 {
		return fmt.Errorf("技能库不可为空")
	}
	Props_Base_Info, err = QueryPropsBaseInfo()
	if err != nil {
		return err
	}
	if len(Props_Base_Info) == 0 {
		return fmt.Errorf("道具库不可为空")
	}
	return InitData()
}

//方便查询，目前返回[]map[string]string内容
func Query(sql string) ([]map[string]string, error) {
	rows, err := dbconn.Query(sql)
	if err != nil {
		return nil, err
	}
	return RowsToMapList(rows)
}

//数据库查询rows转换成maplist类型
func RowsToMapList(rows *sql.Rows) ([]map[string]string, error) {
	var maplist []map[string]string = make([]map[string]string, 0)
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colsize := len(cols)
	for rows.Next() {
		values := make([]sql.NullString, colsize)
		scanArgs := make([]interface{}, colsize)
		for i := 0; i < colsize; i++ {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		oneMap := make(map[string]string)
		for k, v := range values {
			oneMap[cols[k]] = v.String
		}
		maplist = append(maplist, oneMap)
	}
	return maplist, nil
}

func Find(sql string) (map[string]string, error) {
	maplist, err := Query(sql)
	if err != nil {
		log.Println("error:", err.Error(), sql)
		return nil, fmt.Errorf("请稍后")
	}
	if len(maplist) == 0 {
		return nil, nil
	}
	return maplist[0], nil
}

func Count(sql string) (int, error) {
	map_data, err := Find(sql)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(map_data["num"])
	if err != nil {
		return 0, err
	}
	return num, nil
}

func GetHeroByName(name string) map[string]string {
	for _, hero := range Hero_Base_Info {
		if hero["name"] == name {
			return hero
		}
	}
	return nil
}

func GetEqptInfoByName(name string) map[string]string {
	for _, eqpt := range Eqpt_Base_Info {
		if eqpt["name"] == name {
			return eqpt
		}
	}
	return nil
}

func GetSkillInfoByName(name string) map[string]string {
	for _, skill := range Skill_Base_Info {
		if skill["name"] == name {
			return skill
		}
	}
	return nil
}

func GetPropsInfoByName(name, quality string) map[string]string {
	for _, props := range Props_Base_Info {
		if props["name"] == name && props["quality"] == quality {
			return props
		}
	}
	return nil
}
