package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	var db *gorm.DB
	var err error

	if conf.Dev {
		db, err = gorm.Open(sqlite.Open("robot.db"), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(conf.DSN), &gorm.Config{})
	}

	if err != nil {
		log.Println(err)
	}

	if err := db.AutoMigrate(&User{}, &KeyValue{}, &Alpha{}, &Node{}); err != nil {
		panic(err.Error())
	}

	return db
}
