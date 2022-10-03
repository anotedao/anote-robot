package main

import (
	"log"

	"github.com/anonutopia/gowaves"
	macaron "gopkg.in/macaron.v1"
	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

var conf *Config

var bot *telebot.Bot

var anc *gowaves.WavesNodeClient

var db *gorm.DB

var m *macaron.Macaron

var monitor *Monitor

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	m = initMacaron()

	bot = initTelegramBot()

	initCommands()

	anc = initAnote()

	db = initDb()

	monitor = initMonitor()

	initAnoteToday()

	// dataTransaction("3AUcqi5gEWtFyQDy4nLvfm13JFnCBwpGHXT", nil, nil, nil)

	log.Println("AnoteRobot started.")

	bot.Start()
}
