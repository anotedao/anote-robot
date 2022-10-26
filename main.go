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

var pc *PriceClient

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	m = initMacaron()

	bot = initTelegramBot()

	initCommands()

	anc = initAnote()

	db = initDb()

	monitor = initMonitor()

	pc = initPriceClient()

	initAnoteToday()

	// val := int64(Fee * 5)
	// dataTransaction2("%s__3AE23gbkTz3hgvBKgEkpe4cRqcFKgbi2Sns", nil, nil, nil)

	log.Println("AnoteRobot started.")

	bot.Start()
}
