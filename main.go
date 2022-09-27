package main

import (
	"log"

	"github.com/anonutopia/gowaves"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var bot *telebot.Bot

var anc *gowaves.WavesNodeClient

var db *gorm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	bot = initTelegramBot()

	initCommands()

	anc = initAnote()

	db = initDb()

	initMonitor()

	// dataTransaction("3AKGP29V8Pjh5VekzXq1SnwWXjMkQm7Zf9h", nil, nil, nil)

	bot.Start()
}
