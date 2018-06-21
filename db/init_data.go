package db

import (
	"../models"
	"fmt"
	_ "gosql/sqlite3"
	"io/ioutil"
	"log"
	//"strings"
)

//数据过少，自动生成一些玩家数据
func InitData() error {
	/*
		决斗列表用户数过少，生成一些电脑玩家，把决斗列表加到200个，
		这些用户都有六个英雄，
		1-10装备技能随机，
		11-50只有装备无技能，
		51-200无技能无装备
	*/
	log.Println("开始数据初始化")
	duel_count := DuelHeroCount()
	if duel_count < 200 {
		//清空
		DuelDelete()
		for i := 1; i <= 200; i++ {
			//生成用户  (注意：名字以后自动生成)
			name := fmt.Sprintf("大海贼%d", i)
			userId_str := AddComputerPlayer(models.GetRandUser(), "e10cabdc3949ba59abbe56e057f20f883e", name)
			hero_list := make([]map[string]string, 0)
			//生成英雄 生成六个hero
			for len(hero_list) < 6 {
				hero := models.RandGet(getRandQuality(i), Hero_Base_Info)
				//看看是否已经存在
				has := false
				for _, onehero := range hero_list {
					if hero["id"] == onehero["id"] {
						has = true
						break
					}
				}
				//如果不存在则添加
				if !has {
					hero_list = append(hero_list, hero)
				}
				AddHero(userId_str, hero)
			}
			baghero_list, err := BagHeroByUserid(userId_str)
			if err != nil {
				log.Panic(err.Error())
			}
			//生成阵容
			for teamindex := 1; teamindex <= 6; teamindex++ {
				err := SetHero(userId_str, baghero_list[teamindex-1]["id"], fmt.Sprint(teamindex))
				if err != nil {
					log.Panic(err.Error())
				}
			}
			if i <= 50 {
				atk_size, def_size, hp_size := 0, 0, 0
				//生成装备
				for atk_size < 6 || def_size < 6 || hp_size < 6 {
					tmp_eqpt := models.RandGet(getRandQuality(i), Eqpt_Base_Info)
					if tmp_eqpt["typename"] == "武器" && atk_size < 6 {
						err := AddEqpt(userId_str, tmp_eqpt)
						if err != nil {
							log.Panic(err.Error())
						}
						atk_size++
					}
					if tmp_eqpt["typename"] == "防具" && def_size < 6 {
						err := AddEqpt(userId_str, tmp_eqpt)
						if err != nil {
							log.Panic(err.Error())
						}
						def_size++
					}
					if tmp_eqpt["typename"] == "饰品" && hp_size < 6 {
						err := AddEqpt(userId_str, tmp_eqpt)
						if err != nil {
							log.Panic(err.Error())
						}
						hp_size++
					}
				}
				//已经添加到库 ，然后从库中查询出来
				atk_list, err := BagEqptByUserid(userId_str, "武器")
				if err != nil {
					log.Panic(err.Error())
				}
				def_list, err := BagEqptByUserid(userId_str, "防具")
				if err != nil {
					log.Panic(err.Error())
				}
				hp_list, err := BagEqptByUserid(userId_str, "饰品")
				if err != nil {
					log.Panic(err.Error())
				}
				//再设置到阵容
				for k, atk_eqpt := range atk_list {
					SetEqptSkill(userId_str, fmt.Sprint(k+1), "eqpt1", atk_eqpt["id"])
				}
				for k, def_eqpt := range def_list {
					SetEqptSkill(userId_str, fmt.Sprint(k+1), "eqpt2", def_eqpt["id"])
				}
				for k, hp_eqpt := range hp_list {
					SetEqptSkill(userId_str, fmt.Sprint(k+1), "eqpt3", hp_eqpt["id"])
				}
			}
			if i <= 10 {
				//生成技能
				for k := 0; k < 6; k++ {
					skill := models.RandGet(getRandQuality(i), Skill_Base_Info)
					err := AddSkill(userId_str, skill)
					if err != nil {
						log.Panic(err.Error())
					}
				}
				//已经添加到库 ，然后从库中查询出来
				skill_list, err := BagSkillByUserid(userId_str)
				if err != nil {
					log.Panic(err.Error())
				}
				//再设置到阵容
				for k, skill := range skill_list {
					SetEqptSkill(userId_str, fmt.Sprint(k+1), "skill1", skill["id"])
				}
			}
			//添加到决斗中
			err = AddDuel(i, userId_str, name)
			if err != nil {
				log.Panic(err.Error())
			}
		}
	}
	//
	NotFindDoExec("select count(*) as num from config where [key]='newday'",
		"insert into config([key],[value],[time],[des]) values('newday','','','判断是不是新的一天，新的一天，前一天数据要更新掉')")

	//sign data init
	_, err := GetConn().Exec(DEL_SIGN_SQL)
	if err != nil {
		log.Panicln(err.Error())
	}
	bSign, err := ioutil.ReadFile("./initsql/init_sign.sql")
	if err != nil {
		log.Panicln(err.Error())
	}
	_, err = GetConn().Exec(string(bSign))
	if err != nil {
		log.Panicln(err.Error())
	}
	//props初始化
	_, err = GetConn().Exec(DROP_PROPS_SQL)
	bProps, err := ioutil.ReadFile("./table/props.sql")
	if err != nil {
		log.Panicln(err.Error())
	}
	_, err = GetConn().Exec(string(bProps))
	if err != nil {
		log.Panicln(err.Error())
	}
	bProps, err = ioutil.ReadFile("./initsql/init_props.sql")
	if err != nil {
		log.Panicln(err.Error())
	}
	_, err = GetConn().Exec(string(bProps))
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println("数据初始化完成")
	return nil
}

//生成品质等级
func getRandQuality(index int) string {
	if index <= 10 {
		return "A"
	}
	if index <= 50 {
		return "AB"
	}
	return "ABC"
}

func NotFindDoExec(query, exec string) {
	num, err := Count(query)
	if err != nil {
		log.Panic(err.Error())
	}
	if num != 1 {
		_, err = GetConn().Exec(exec)
		if err != nil {
			log.Panic(err.Error())
		}
	}
}
