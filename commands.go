package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func initCommands() {
	bot.Handle("/hello", helloCommand)
}

func helloCommand(m *tb.Message) {
	prices := fmt.Sprintf("%#v", pc.Prices)
	bot.Send(m.Sender, prices)
}
