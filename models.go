package main

import (
	"time"

	"gorm.io/gorm"
)

type Miner struct {
	gorm.Model
	Address          string `gorm:"size:255;uniqueIndex"`
	LastNotification time.Time
}
