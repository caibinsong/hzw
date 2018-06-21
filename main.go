package main

import (
	"./controllers"
	"./db"
	"./run"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

/**
服务端第一版，以快速开发为主
*/
var ControllerList []interface{}

//name是有包路径在的，从name中提取最后一个.后面的类名
func GetClassName(name string) string {
	lastindex := strings.LastIndex(name, ".")
	if lastindex >= 0 {
		return name[lastindex+1:]
	}
	return name
}

//添加自动路由
func AddAutoRouter(controller interface{}) {
	//获取类型名
	datatype := GetClassName(reflect.TypeOf(controller).String())
	//判断是否类名最后是Controller结尾的
	if strings.HasSuffix(datatype, "Controller") {
		ControllerList = append(ControllerList, controller)
		//获取Controller前面的字符串来当位上一级目录
		controllerName := strings.ToLower(datatype[:len(datatype)-len("Controller")])

		reflectVal := reflect.ValueOf(controller)
		//遍历类中的方法
		for i := 0; i < reflectVal.NumMethod(); i++ {
			funPtr := reflectVal.Method(i)
			//如果方法是func(http.ResponseWriter, *http.Request)类型的就认为要添加路由
			if funPtr.Type().String() == "func(http.ResponseWriter, *http.Request)" {
				http.HandleFunc("/"+controllerName+"/"+strings.ToLower(reflectVal.Type().Method(i).Name),
					func(w http.ResponseWriter, req *http.Request) {

						funPtr.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(req)})
					})
			}
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	by, err := ioutil.ReadFile("./conf/app.conf")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("系统启动:" + string(by))
	//添加自动路由
	AddAutoRouter(&controllers.UserController{})     //用户
	AddAutoRouter(&controllers.MoneyController{})    //钱
	AddAutoRouter(&controllers.StrengthController{}) //体力
	AddAutoRouter(&controllers.RaffleController{})   //抽奖
	AddAutoRouter(&controllers.TeamController{})     //阵容
	AddAutoRouter(&controllers.BagController{})      //背包
	AddAutoRouter(&controllers.DuelController{})     //决斗
	AddAutoRouter(&controllers.SignController{})     //签到
	AddAutoRouter(&controllers.PropsController{})    //道具

	//数据库初始化
	err = db.DBConn()
	if err != nil {
		log.Println(err.Error())
		return
	}
	go run.NewDay()
	err = http.ListenAndServe(":"+string(by), nil)
	if err != nil {
		log.Println("服务启动失败:", err)
		return
	}
	log.Println("系统退出")
}
