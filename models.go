package main

import "gorm.io/gorm"

// User represents Telegram user
type User struct {
	gorm.Model
	Address         string `gorm:"size:255;uniqueIndex"`
	NotificationDay uint64
}
