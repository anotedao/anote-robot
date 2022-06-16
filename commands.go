package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func initCommands() {
	bot.Handle("/hello", helloCommand)
}

func helloCommand(m *tb.Message) {
	prices := fmt.Sprintf("Well, hello %s!", m.Sender.FirstName)
	bot.Send(m.Chat, prices)
}
