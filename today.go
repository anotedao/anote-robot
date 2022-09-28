package main

import (
	"fmt"
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
		at.sendAd(fmt.Sprintf(defaultAd, 178))

		time.Sleep(time.Second * MonitorTick)
	}
}

func initAnoteToday() {
	at := &AnoteToday{}
	go at.start()
}

var defaultAd = `<b><u>🔴  ANOTE 2.0 IS LIVE</u></b>    🚀

We are proud to announce that Anote 2.0 is now available for mining.

We now have our own wallet (anote.one) which is used both as a wallet and a tool for mining. Stay tuned for more exciting news!

________________________
Daily Mining Code: %d
`
