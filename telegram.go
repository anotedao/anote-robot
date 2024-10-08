package main

import (
	"log"
	"net/url"
	"strings"
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

func initTelegramBot2() *telebot.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:     conf.TelegramAPIKey2,
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
	message = "anote-robot:" + getCallerInfo() + message
	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}
	bot.Send(rec, message)
}

func logTelegramSilent(message string) {
	message = "anote-robot:" + getCallerInfo() + message
	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}
	bot.Send(rec, message, telebot.Silent)
}

func notificationTelegram(message string) {
	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}
	bot2.Send(rec, message, telebot.NoPreview)
}

func notificationTelegramPin(message string) {
	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}
	m, _ := bot2.Send(rec, message)
	bot2.Pin(m)
}

func notificationTelegramTeam(message string) {
	rec := &telebot.Chat{
		ID: int64(TelGroup),
	}
	bot2.Send(rec, message, telebot.NoPreview)
}

func notificationTelegramTeamPin(message string) {
	rec := &telebot.Chat{
		ID: int64(TelGroup),
	}
	m, _ := bot2.Send(rec, message)
	bot2.Pin(m)
}

func notificationTelegramGroup(message string) {
	rec := &telebot.Chat{
		ID: int64(TelAnon),
	}
	bot2.Send(rec, message, telebot.NoPreview)
}

func notificationTelegramGroupBalkan(message string) {
	rec := &telebot.Chat{
		ID: int64(TelBalkan),
	}
	bot2.Send(rec, message, telebot.NoPreview, telebot.Silent)
}

func notificationTelegramGroupBalkanPin(message string) {
	rec := &telebot.Chat{
		ID: int64(TelBalkan),
	}
	m, _ := bot2.Send(rec, message)
	bot2.Pin(m, telebot.Silent)
}

func notificationTelegramGroupPin(message string) {
	rec := &telebot.Chat{
		ID: int64(TelAnon),
	}
	m, _ := bot2.Send(rec, message)
	bot2.Pin(m)
}

func logTelegramService(message string) error {
	m, _ := url.QueryUnescape(message)
	message, _ = url.PathUnescape(m)
	var err error

	rec := &telebot.Chat{
		ID: int64(TelAnonOps),
	}

	if strings.Contains(message, "no data for this key") {
		_, err = bot.Send(rec, message, telebot.NoPreview, telebot.Silent)
	} else {
		_, err = bot.Send(rec, message)
	}
	return err
}
