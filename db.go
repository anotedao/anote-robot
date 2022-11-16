package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDb() *gorm.DB {
	var db *gorm.DB
	var err error
	dbconf := gorm.Config{}

	dbconf.Logger = logger.Default.LogMode(logger.Error)

	db, err = gorm.Open(sqlite.Open("robot.db"), &dbconf)

	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if err := db.AutoMigrate(&Miner{}, &User{}, &KeyValue{}); err != nil {
		panic(err.Error())
	}

	return db
}
