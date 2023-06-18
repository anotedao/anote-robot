package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/dustin/go-humanize"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

func initCommands() {
	bot.Handle("/help", helpCommand)
	bot.Handle("/start", startCommand)
	bot.Handle("/stats", statsCommand)
	bot.Handle("/miner", myStatsCommand)
	bot.Handle("/code", codeCommand)
	bot.Handle("/bo", batteryCommand)
	bot.Handle("/ref", refCommand)
	bot.Handle("/withdraw", withdrawCommand)
	bot.Handle(telebot.OnUserJoined, userJoined)
	bot.Handle(telebot.OnText, mineCommand)
}

func helpCommand(c telebot.Context) error {
	adnum, err := getData2("%s__adnum", nil)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	saveUser(c)
	m := c.Message()

	help := fmt.Sprintf(`⭕️ <b><u>Anote Mining Tutorial</u></b>
	
	To start mining Anote, follow these simple steps:

	  - read the daily mining code from <a href="https://t.me/AnoteAds/%d">AnoteAds</a> channel
	  - open @AnoteRobot and click start if you already haven't
	  - send the daily mining code to AnoteRobot as a message
	  
	And that's it, you are now mining Anote. 🚀`, adnum.(int64))

	_, err = bot.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func startCommand(c telebot.Context) error {
	saveUser(c)
	m := c.Message()
	split := strings.Split(m.Text, " ")
	response := ""
	tid := strconv.Itoa(int(m.Chat.ID))

	adnum, err := getData2("%s__adnum", nil)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if len(split) == 2 {
		if saveTelegram(split[1], tid) != 0 {
			response = "There's already an Anote address attached to this Telegram account."
		} else {
			saveTelegram("none", tid)
			response = fmt.Sprintf(`⭕️ <b><u>Welcome to Anote!</u></b> 🚀
			
Start mining by reading the daily mining code in <a href="https://t.me/AnoteAds/%d">AnoteAds</a> channel and sending it back here to activate the mining cycle.
		
Join @AnoteDAO group for help and support.`,
				adnum)
		}
	} else {
		saveTelegram("none", tid)
		response = fmt.Sprintf(`⭕️ <b><u>Welcome to Anote!</u></b> 🚀
		
Start mining by reading the daily mining code in <a href="https://t.me/AnoteAds/%d">AnoteAds</a> channel and sending it back here to activate the mining cycle.
		
Join @AnoteDAO group for help and support.`,
			adnum)
	}

	_, err = bot.Send(m.Chat, response, telebot.NoPreview)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	return nil
}

func statsCommand(c telebot.Context) error {
	saveUser(c)
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

	stats := getStats()

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := proto.MustAddressFromString(MobileAddress)
	addrt := proto.MustAddressFromString(TelegramAddress)

	total, _, err := cl.Addresses.Balance(ctx, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	totalt, _, err := cl.Addresses.Balance(ctx, addrt)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	log.Println(prettyPrint(stats))

	basicAmount := float64(0)
	basicAmountT := float64(0)

	if stats.ActiveUnits > 0 {
		basicAmount = float64((float64(total.Balance) / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4)))) / MULTI8
		basicAmountT = float64((float64(totalt.Balance) / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4)))) / MULTI8
	} else {
		basicAmount = float64((float64(total.Balance)) / MULTI8)
		basicAmountT = float64((float64(totalt.Balance)) / MULTI8)
	}

	log.Println(basicAmount)
	log.Println(basicAmountT)

	s := fmt.Sprintf(`⭕️ <u><b>Anote Basic Stats</b></u>
	
	<b>Active Miners:</b> %d
	<b>Holders:</b> %d
	<b>Price:</b> $%.2f
	<b>Telegram Amount:</b> %.8f
	<b>Mobile Amount:</b> %.8f
	
	<b>Mined:</b> %s ANOTE
	<b>Community:</b> %s ANOTE
	<b>In Circulation:</b> %s ANOTE
	
	<b>Referred Miners:</b> %d
	<b>Payout Miners:</b> %d
	<b>Inactive Miners:</b> %d`,
		stats.ActiveMiners,
		stats.Holders,
		pc.AnotePrice,
		basicAmountT,
		basicAmount,
		humanize.Comma(mined),
		humanize.Comma(int64(balance)),
		humanize.Comma(circulation),
		stats.ActiveReferred,
		stats.PayoutMiners,
		stats.InactiveMiners)

	bot.Send(m.Chat, s)

	return nil
}

func userJoined(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	msg := fmt.Sprintf("Hello, %s! Welcome to Anote community! 🚀\n\nHere are some resources to get you started:\n\nAnote Wallet: anote.one\nBlockchain Explorer: explorer.anotedao.com\nWebsite: anotedao.com\nMining Tutorial: anotedao.com/mining\nRun a Node: anotedao.com/node\n\n<u>Other Anote Communities:</u>\n\n@AnoteBalkan\n@AnoteAfrica\n@AnoteChina", m.Sender.FirstName)

	bot.Send(m.Chat, msg, telebot.NoPreview)

	return nil
}

func deleteCommand(c telebot.Context) error {
	saveUser(c)
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
		bot.Send(msg.Chat, "Please send this command as a direct message to @AnoteRobot.")
	}

	return nil
}

func codeCommand(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	adnum, err := getData2("%s__adnum", nil)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	help := fmt.Sprintf("<a href=\"https://t.me/AnoteAds/%d\"><strong><u>Click here</u></strong></a>, daily mining code is at the bottom of the last announcement.", adnum.(int64))

	_, err = bot.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func saveUser(c telebot.Context) {
	m := c.Message()
	user := &User{}
	db.FirstOrCreate(user, User{TelegramId: m.Chat.ID})
	user.TelUsername = m.Sender.Username
	user.TelName = m.Sender.FirstName + " " + m.Sender.LastName
	user.TelDump = prettyPrint(m.Sender)
	db.Save(user)
}

func batteryCommand(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	help := "To achieve 100% AINT Miner health and receive full amount of anotes, disable battery optimization on AINT Miner. You can learn how to do that here:\n\nanotedao.com/battery"

	_, err := bot.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func mineCommand(c telebot.Context) error {
	var err error
	if c.Message().Private() {
		message := telegramMine(c.Message().Text, c.Chat().ID)
		_, err = bot.Send(c.Chat(), message, telebot.NoPreview)
	}

	return err
}

func myStatsCommand(c telebot.Context) error {
	msg := c.Message()
	var err error
	abr := &gowaves.AddressesBalanceResponse{}

	miner := getMiner(c.Message().Sender.ID)
	if strings.HasPrefix(miner.Address, "3A") {
		abr, err = anc.AddressesBalance(miner.Address)
		if err != nil {
			log.Println(err.Error())
			logTelegram(err.Error())
		}
	}

	blocks := 1410 - int64(miner.Height) + miner.MiningHeight
	if blocks < 0 {
		blocks = 0
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%dm", blocks))
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	message := fmt.Sprintf(`⭕️ <b><u>Your Anote Stats</u></b>

	<b>Mined Telegram:</b> %.8f ANOTE
	<b>Mined Mobile:</b> N/A
	<b>Address Balance:</b> %.8f ANOTE

	<b>Cycle Blocks Left:</b> %d
	<b>Cycle Time Left:</b> %02d:%02d

	<b>Referred Total:</b> %d
	<b>Referred Active:</b> %d
	<b>Referred Confirmed:</b> %d
	
	<b><u>Other Commands</u></b>

	/ref - Your Anote referral link
	/withdraw - Withdraw your mined anotes`,
		float64(miner.MinedTelegram)/float64(MULTI8),
		float64(abr.Balance)/float64(MULTI8),
		blocks,
		int(duration.Hours()),
		int(duration.Minutes())%60,
		miner.Referred,
		miner.Active,
		miner.Confirmed)

	if !msg.Private() {
		message = "Please send this command as a direct message to @AnoteRobot."
	}

	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func refCommand(c telebot.Context) error {
	msg := c.Message()
	var err error

	if !msg.Private() {
		message := "Please send this command as a direct message to @AnoteRobot."
		_, err = bot.Send(c.Chat(), message, telebot.NoPreview)
		return err
	}

	miner := getMiner(c.Message().Sender.ID)

	message := fmt.Sprint("Your Anote referral link:")
	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)

	message = fmt.Sprintf("https://t.me/AnoteRobot?start=%d", miner.ID)
	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func withdrawCommand(c telebot.Context) error {
	// msg := c.Message()
	var err error
	// message := ""

	// if !msg.Private() {
	// 	message := "Please send this command as a direct message to @AnoteRobot."
	// 	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)
	// 	return err
	// }

	// miner := getMiner(c.Message().Sender.ID)

	// if miner.MinedMobile+miner.MinedTelegram > Fee {
	// 	if strconv.Itoa(int(miner.TelegramId)) == miner.Address {
	// 		message = "To withdraw your funds, please open account on anote.one wallet and connect it to the bot!"
	// 	} else {
	// 		message = "Your funds have been sent to your address. 🚀"
	// 	}
	// } else {
	// 	message = "You don't have enough funds to withdraw. Please try later!"
	// }

	message := "This command is under construction."
	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)

	return err
}
