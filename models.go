package main

import (
	"time"

	"gorm.io/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;unique_index"`
	ValueInt uint64
	ValueStr string
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address         string `sql:"size:255;unique_index"`
	TelegramID      int    `sql:"unique_index"`
	ReferralID      uint
	Referral        *User
	MiningActivated *time.Time
	LastStatus      *time.Time
	MinedAnotes     int
	Mining          bool `sql:"DEFAULT:false"`
	LastWithdraw    *time.Time
	Language        string `sql:"size:255;"`
	MiningWarning   *time.Time
	Nickname        string `sql:"size:255;unique_index"`
	FbPostLink      string `sql:"size:255;"`
	SentAint        bool   `sql:"DEFAULT:false"`
	LastFbQuest     *time.Time
	SentFbAnotes    bool   `sql:"DEFAULT:false"`
	TwPostLink      string `sql:"size:255;"`
	LastTwQuest     *time.Time
	SentTwAnotes    bool `sql:"DEFAULT:false"`
	UpdatedAddress  bool `sql:"DEFAULT:false"`
}
