package main

import (
	"log"

	"github.com/anonutopia/gowaves"
	"gopkg.in/tucnak/telebot.v2"
)

var conf *Config

var bot *telebot.Bot

var anc *gowaves.WavesNodeClient

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	bot = initTelegramBot()

	initCommands()

	anc = initAnote()

	bot.Start()
}
