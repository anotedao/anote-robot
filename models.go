package main

import "gorm.io/gorm"

type Miner struct {
	gorm.Model
	Address         string `gorm:"size:255;uniqueIndex"`
	NotificationDay uint64
}
