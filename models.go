package main

import (
	"gorm.io/gorm"
)

// type Miner struct {
// 	gorm.Model
// 	Address          string `gorm:"size:255;uniqueIndex"`
// 	LastNotification time.Time
// }

type User struct {
	gorm.Model
	TelegramId  int64
	TelUsername string `gorm:"size:255"`
	TelName     string `gorm:"size:255"`
	TelDump     string `gorm:"size:512"`
}

type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt uint64 `gorm:"type:int"`
	ValueStr string `gorm:"type:string"`
}

type Alpha struct {
	gorm.Model
	Address string `gorm:"size:255;uniqueIndex"`
}
