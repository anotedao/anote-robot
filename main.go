package main

import (
	"github.com/anonutopia/gowaves"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var bot *telebot.Bot

var db *gorm.DB

var wnc *gowaves.WavesNodeClient

var wmc *gowaves.WavesMatcherClient

var pc *PriceClient

func main() {
	conf = initConfig()

	bot = initTelegramBot()

	db = initDb()

	wnc, wmc = initWaves()

	pc = initPriceClient()

	initCommands()

	bot.Start()
}
