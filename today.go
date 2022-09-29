package main

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/telebot.v3"
)

type AnoteToday struct {
}

func (at *AnoteToday) sendAd(ad string) {
	var channeId int64
	if conf.Dev {
		channeId = TelDevAnoteToday
	} else {
		channeId = TelAnoteToday
	}
	r := &telebot.Chat{
		ID: channeId,
	}

	bot.Send(r, ad)
}

func (at *AnoteToday) start() {
	for {
		if at.isNewCycle() {
			code := at.generateNewCode()

			ad := at.getAd()

			at.sendAd(fmt.Sprintf(ad, code))
		}

		time.Sleep(time.Second * MonitorTick)
	}
}

func (at *AnoteToday) isNewCycle() bool {
	ks := &KeyValue{Key: "lastAdDay"}
	db.FirstOrCreate(ks, ks)
	today := time.Now().Day()

	if ks.ValueInt != uint64(today) && time.Now().Hour() == SendAdHour {
		ks.ValueInt = uint64(today)
		db.Save(ks)

		return true
	}

	return false
}

func (at *AnoteToday) generateNewCode() int {
	ks := &KeyValue{Key: "dailyCode"}
	db.FirstOrCreate(ks, ks)

	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 999

	code := rand.Intn(max-min+1) + min

	ks.ValueInt = uint64(code)
	db.Save(ks)

	return code
}

func (at *AnoteToday) getAd() string {
	return defaultAd
}

func initAnoteToday() {
	at := &AnoteToday{}
	go at.start()
}

var defaultAd = `<b><u>🔴  ANOTE 2.0 IS NOW LIVE</u></b>    🚀

We are proud to announce that Anote 2.0 is now available for mining.

We now have our own wallet (anote.one) which is used both as a wallet and a tool for mining. Stay tuned for more exciting news, information and tutorials!

________________________
Daily Mining Code: %d
`
