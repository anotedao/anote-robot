package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"gopkg.in/telebot.v3"
)

func initCommands() {
	bot.Handle("/hello", helloCommand)
	bot.Handle("/start", startCommand)
	bot.Handle("/stats", statsCommand)
}

func helloCommand(c telebot.Context) error {
	m := c.Message()
	hello := fmt.Sprintf("Well, hello %s!", m.Sender.FirstName)
	bot.Send(m.Chat, hello)
	return nil
}

func startCommand(c telebot.Context) error {
	m := c.Message()
	split := strings.Split(m.Text, " ")
	response := ""

	if len(split) == 2 {
		telId := strconv.Itoa(int(m.Chat.ID))
		encTelegram := EncryptMessage(telId)
		err := dataTransaction(split[1], &encTelegram, nil, nil)
		if err != nil {
			log.Println(err)
		}
		response = "You have successfully connected your anote.one wallet to the bot. 🚀 Please restart or reload the wallet to start mining!\n\nYou can find daily mining code in @AnoteToday channel."
	} else {
		response = "Please run this bot from anote.one wallet (click the blue briefcase icon and then connect Telegram)!"
	}

	bot.Send(m.Chat, response)

	return nil
}

func statsCommand(c telebot.Context) error {
	m := c.Message()
	bh, err := anc.BlocksHeight()
	if err != nil {
		log.Println(err.Error())
	}
	mined := int64(bh.Height + 1000)

	abr, err := anc.AddressesBalance(COMMUNITY_ADDR)
	if err != nil {
		log.Println(err.Error())
	}

	abr2, err := anc.AddressesBalance(GATEWAY_ADDR)
	if err != nil {
		log.Println(err.Error())
	}

	balance := (abr.Balance / int(SATINBTC)) + (abr2.Balance / int(SATINBTC))
	circulation := mined - int64(balance)

	stats := fmt.Sprintf(
		"<u><b>🚀 Anote Basic Stats</b></u>\n\nMined: %s ANOTE\nCommunity: %s ANOTE\nIn Circulation: %s ANOTE",
		humanize.Comma(mined),
		humanize.Comma(int64(balance)),
		humanize.Comma(circulation))

	bot.Send(m.Chat, stats)

	return nil
}
