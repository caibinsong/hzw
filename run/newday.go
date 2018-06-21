package run

import (
	"../db"
	"log"
	"time"
)

func NewDay() {
	for true {
		if db.IsNewDay() {
			log.Println("is newday")
			db.AdventureNewDay()
		}
		time.Sleep(60 * time.Second)
	}
}
