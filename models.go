package main

import (
	"time"

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

type Miner struct {
	gorm.Model
	Address                string `gorm:"size:255;uniqueIndex"`
	LastNotification       time.Time
	LastNotificationWeekly time.Time `gorm:"default:'2023-06-17 23:00:00.797487649+00:00'"`
	TelegramId             int64     `gorm:"uniqueIndex"`
	MiningHeight           int64
	MiningTime             time.Time
	ReferralID             uint `gorm:"index"`
	Balance                uint64
	MinedTelegram          uint64
	MinedMobile            uint64
	LastPing               time.Time
	PingCount              int64
	IpAddresses            []*IpAddress `gorm:"many2many:miner_ip_addresses;"`
	UpdatedApp             bool         `gorm:"default:false"`
	LastInvite             time.Time
	BatteryNotification    bool `gorm:"default:false"`
	Cycles                 uint64
}

type IpAddress struct {
	gorm.Model
	Address string   `gorm:"size:255;uniqueIndex"`
	Miners  []*Miner `gorm:"many2many:miner_ip_addresses;"`
}
