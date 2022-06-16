package main

import (
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initTelegramBot() *telebot.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:   conf.TelegramAPIKey,
		Poller:  &tb.LongPoller{Timeout: T_POLLER_TIMEOUT * time.Second},
		Verbose: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	return b
}
