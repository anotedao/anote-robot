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

	// encId, _ := getData("3ANmnLHt8mR9c36mdfQVpBtxUs8z1mMAHQW")
	// telId := DecryptMessage(encId.(string))
	// idNum, _ := strconv.Atoi(telId)
	// // rec := telebot.Recipient
	// rec := &telebot.Chat{
	// 	ID: int64(idNum),
	// }

	// bot.Send(rec, "hello world")

	// dataTransaction("3AKGP29V8Pjh5VekzXq1SnwWXjMkQm7Zf9h", nil, nil, nil)

	bot.Start()
}
