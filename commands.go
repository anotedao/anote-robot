package main

import (
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initCommands() {
	bot.Handle("/hello", helloCommand)
	bot.Handle("/stats", statsCommand)
}

func helloCommand(m *tb.Message) {
	hello := fmt.Sprintf("Well, hello %s!", m.Sender.FirstName)
	bot.Send(m.Chat, hello)
}

func statsCommand(m *tb.Message) {
	bh, err := anc.BlocksHeight()
	if err != nil {
		log.Println(err.Error())
	}
	mined := int64(bh.Height + 1000)

	abr, err := anc.AddressesBalance(COMMUNITY_ADDR)
	if err != nil {
		log.Println(err.Error())
	}
	balance := abr.Balance / int(SATINBTC)
	circulation := mined - int64(balance)

	stats := fmt.Sprintf("<u><b>🚀 Anote Basic Stats</b></u>\n\nMined: %s ANOTE\nCommunity: %s ANOTE\nIn Circulation: %s ANOTE", humanize.Comma(mined), humanize.Comma(int64(balance)), humanize.Comma(circulation))
	bot.Send(m.Chat, stats)
}
