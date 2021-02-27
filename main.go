package main

import (
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var b *telebot.Bot

var db *gorm.DB

func main() {
	conf = initConfig()

	b = initTelegramBot()

	db = initDb()

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, string(m.Text))
	})

	b.Start()
}
