package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

func initCommands() {
	bot.Handle("/help", helpCommand)
	bot.Handle("/start", startCommand)
	bot.Handle("/stats", statsCommand)
	bot.Handle("/delete", deleteCommand)
	bot.Handle(telebot.OnUserJoined, userJoined)
}

func helpCommand(c telebot.Context) error {
	m := c.Message()

	help := `⭕️ <b><u>Anote Mining Tutorial</u></b>
	
	To start mining Anote, follow these simple steps:

	  - read daily mining code on @AnoteToday Telegram channel
	  - open anote.one wallet
	  - click blue briefcase icon
	  - connect Telegram bot by clicking the button
	  - enter daily mining code and captcha
	  - click mine button
	  - join @AnoteDigital group for help and support
	  
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
		if monitor.minerExists(m.Chat.ID) {
			log.Println(telId)
			response = "There's already an Anote address attached to this Telegram account."
		} else {
			encTelegram := EncryptMessage(telId)
			minerData := "%s%d%s" + Sep + encTelegram
			err := dataTransaction(split[1], &minerData, nil, nil)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}
			response = "You have successfully connected your anote.one wallet to the bot. 🚀 Please restart or reload the wallet to start mining!\n\nYou can find daily mining code in @AnoteToday channel.\n\nJoin @AnoteDigital group for help and support."
		}
	} else {
		response = "Please run this bot from anote.one wallet (click the blue briefcase icon and then connect Telegram)!\n\nJoin @AnoteDigital group for help and support."
	}

	bot.Send(m.Chat, response)

	go monitor.loadMiners()

	return nil
}

func statsCommand(c telebot.Context) error {
	m := c.Message()
	bh, err := anc.BlocksHeight()
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}
	mined := int64(bh.Height + 1000)

	abr, err := anc.AddressesBalance(COMMUNITY_ADDR)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	abr2, err := anc.AddressesBalance(GATEWAY_ADDR)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	abr3, err := anc.AddressesBalance(MobileAddress)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	balance := (abr.Balance / int(SATINBTC)) + (abr2.Balance / int(SATINBTC)) + (abr3.Balance / int(SATINBTC))
	circulation := mined - int64(balance)

	miner := getMiner(conf.TelegramAPIKey)

	stats := fmt.Sprintf(
		"⭕️ <u><b>Anote Basic Stats</b></u>\n\nMined: %s ANOTE\nCommunity: %s ANOTE\nIn Circulation: %s ANOTE\nActive Miners: %d\nReferred Miners: %d\nPrice: $%.2f",
		humanize.Comma(mined),
		humanize.Comma(int64(balance)),
		humanize.Comma(circulation),
		miner.ActiveMiners,
		miner.MinRefCount-miner.ActiveMiners,
		pc.AnotePrice)

	bot.Send(m.Chat, stats)

	return nil
}

func userJoined(c telebot.Context) error {
	m := c.Message()

	msg := fmt.Sprintf("Hello, %s! Welcome to Anote community! 🚀\n\nHere are some resources to get you started:\n\nAnote Wallet: anote.one\nBlockchain Explorer: anote.live\nWebsite: anote.digital\nMining Tutorial: anote.digital/mine\nRun a Node: anote.digital/node\n\nIf you are from the Balkans (and you know the language), you can also join our local @AnoteBalkan group.", m.Sender.FirstName)

	bot.Send(m.Chat, msg, telebot.NoPreview)

	return nil
}

func deleteCommand(c telebot.Context) error {
	msg := c.Message()

	if msg.Private() {
		cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		addr, err := proto.NewAddressFromString(MobileAddress)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			return err
		}

		entries, _, err := cl.Addresses.AddressesData(ctx, addr)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			return err
		}

		for _, m := range entries {
			encTel := parseItem(m.ToProtobuf().GetStringValue(), 0)
			telIdStr := DecryptMessage(encTel.(string))
			telId, err := strconv.Atoi(telIdStr)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
				return err
			}

			if telId == int(msg.Chat.ID) {
				err := dataTransaction(m.GetKey(), nil, nil, nil)
				if err != nil {
					log.Println(err)
					logTelegram(err.Error())
					return err
				}
			}
		}

		bot.Send(msg.Chat, "Your account has been successfully disconnected.")
	} else {
		bot.Send(msg.Chat, "Please send this command as a private message to bot.")
	}

	return nil
}
