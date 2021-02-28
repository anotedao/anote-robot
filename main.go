package main

import (
	"github.com/anonutopia/gowaves"
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var b *telebot.Bot

var db *gorm.DB

var wnc *gowaves.WavesNodeClient

var wmc *gowaves.WavesMatcherClient

var pc *PriceClient

func main() {
	conf = initConfig()

	b = initTelegramBot()

	db = initDb()

	wnc, wmc = initWaves()

	pc = initPriceClient()

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, string(m.Text))
	})

	b.Start()
}
