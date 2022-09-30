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
	bot.Handle("/help", helpCommand)
	bot.Handle("/start", startCommand)
	bot.Handle("/stats", statsCommand)
}

func helpCommand(c telebot.Context) error {
	m := c.Message()

	help := `⭕️ <b><u>Anote Mining Tutorial</u></b>
	
	To start mining Anote, follow these simple steps:

	  - read daily mining code on @AnoteToday Telegram channel
	  - open anote.one wallet
	  - click blue briefcase icon
	  - enter daily mining code and captcha
	  - click mine button
	  
	And that's it, you are now mining Anote. 🚀`

	_, err := bot.Send(m.Chat, help)

	return err
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

	abr3, err := anc.AddressesBalance(MobileAddress)
	if err != nil {
		log.Println(err.Error())
	}

	abr4, err := anc.AddressesBalance(TelegramAddress)
	if err != nil {
		log.Println(err.Error())
	}

	balance := (abr.Balance / int(SATINBTC)) + (abr2.Balance / int(SATINBTC)) + (abr3.Balance / int(SATINBTC)) + (abr4.Balance / int(SATINBTC))
	circulation := mined - int64(balance)

	stats := fmt.Sprintf(
		"⭕️ <u><b>Anote Basic Stats</b></u>\n\nMined: %s ANOTE\nCommunity: %s ANOTE\nIn Circulation: %s ANOTE",
		humanize.Comma(mined),
		humanize.Comma(int64(balance)),
		humanize.Comma(circulation))

	bot.Send(m.Chat, stats)

	return nil
}
