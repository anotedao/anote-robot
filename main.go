package main

import (
	"log"

	"github.com/anonutopia/gowaves"
	"gopkg.in/telebot.v3"
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

	initAnoteToday()

	log.Println("AnoteRobot started.")

	bot.Start()
}
