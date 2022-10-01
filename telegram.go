package main

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func initTelegramBot() *telebot.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:     conf.TelegramAPIKey,
		Poller:    &telebot.LongPoller{Timeout: T_POLLER_TIMEOUT * time.Second},
		Verbose:   false,
		ParseMode: "html",
	})

	if err != nil {
		log.Fatal(err)
	}

	return b
}

func logTelegram(message string) {
	message = getCallerInfo() + message
	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}
	bot.Send(rec, message)
}

func logTelegramService(message string) error {
	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}
	_, err := bot.Send(rec, message)
	return err
}
