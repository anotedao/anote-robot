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
	bot2.Handle("/help", helpCommand2)
	bot2.Handle("/stats", statsCommand)
	bot2.Handle("/code", codeCommand)
	bot2.Handle("/bo", batteryCommand)
	bot2.Handle("/alpha", alphaCommand)
	bot2.Handle("/check", checkCommand)
	bot2.Handle(telebot.OnUserJoined, userJoined)
	bot2.Handle(telebot.OnText, checkUserCommand)

	bot.Handle("/start", startCommand)
	bot.Handle("/miner", myStatsCommand)
	bot.Handle("/withdraw", withdrawCommand)
	bot.Handle("/ref", refCommand)
	bot.Handle("/help", helpCommand)
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
	  - join @AnoteDAO group for help and support
	  
	And that's it, you are now mining Anote. 🚀`, adnum.(int64))

	_, err = bot.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func helpCommand2(c telebot.Context) error {
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
	  - join @AnoteDAO group for help and support
	  
	And that's it, you are now mining Anote. 🚀`, adnum.(int64))

	_, err = bot2.Send(m.Chat, help, telebot.NoPreview)

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
			response = "This address is already used."
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
	log.Println(prettyPrint(m))
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
	// addrt := proto.MustAddressFromString(TelegramAddress)

	total, _, err := cl.Addresses.Balance(ctx, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	// totalt, _, err := cl.Addresses.Balance(ctx, addrt)
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	// log.Println(prettyPrint(stats))

	basicAmount := float64(0)
	basicAmountT := float64(0)

	if stats.ActiveUnits > 0 {
		basicAmount = float64((float64(total.Balance) / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4)))) / MULTI8
		basicAmountT = float64((7.2 / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4))))
	} else {
		basicAmount = float64((float64(total.Balance)) / MULTI8)
		basicAmountT = 7.2
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
	<b>In Circulation:</b> %s ANOTE`,
		stats.ActiveMiners,
		stats.Holders,
		pc.AnotePrice,
		basicAmountT,
		basicAmount,
		humanize.Comma(mined),
		humanize.Comma(int64(balance)),
		humanize.Comma(circulation))

	if m.Private() {
		s += fmt.Sprintf(`
	
	<b>Referred Miners:</b> %d
	<b>Payout Miners:</b> %d
	<b>Inactive Miners:</b> %d`,
			stats.ActiveReferred,
			stats.PayoutMiners,
			stats.InactiveMiners)
	}

	bot2.Send(m.Chat, s)

	return nil
}

func userJoined(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	msg := fmt.Sprintf("Hello, %s! Welcome to Anote community! 🚀\n\nHere are some resources to get you started:\n\nAnote Wallet: app.anotedao.com\nBlockchain Explorer: explorer.anotedao.com\nWebsite: anotedao.com\nMining Tutorial: anotedao.com/mine\nRun a Node: anotedao.com/node\n\n<u>Other Anote Communities:</u>\n\n@AnoteBalkan\n@AnoteAfrica\n@AnoteChina", m.Sender.FirstName)

	bot2.Send(m.Chat, msg, telebot.NoPreview)

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

	_, err = bot2.Send(m.Chat, help, telebot.NoPreview)

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

	_, err := bot2.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func mineCommand(c telebot.Context) error {
	var err error
	if c.Message().Private() {
		message := ""
		if strings.HasPrefix(c.Message().Text, "3A") {
			if saveTelegram(c.Message().Text, strconv.Itoa(int(c.Chat().ID))) != 0 {
				message = "This address is already used."
			} else {
				message = "You have successfully connected your Anote wallet. 🚀"
			}
		} else if c.Message().IsForwarded() {
			message = "Forwarded."
		} else {
			message = telegramMine(c.Message().Text, c.Chat().ID)
		}
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

	blocks := 1411 - int64(miner.Height) + miner.MiningHeight
	if blocks < 0 {
		blocks = 0
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%dm", blocks))
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	address := "N/A"
	if strings.HasPrefix(miner.Address, "3A") {
		address = miner.Address
	}

	message := fmt.Sprintf(`⭕️ <b><u>Your Anote Stats</u></b>

	<b>Mined Telegram:</b> %.8f ANOTE
	<b>Mined Mobile:</b> N/A
	<b>Address Balance:</b> %.8f ANOTE

	<b>Cycle Blocks Left:</b> %d
	<b>Cycle Time Left:</b> %02d:%02d

	<b>Referred Total:</b> %d
	<b>Referred Active:</b> %d

	<b><u>Address</u></b>

	%s`,
		float64(miner.MinedTelegram)/float64(MULTI8),
		float64(abr.Balance)/float64(MULTI8),
		blocks,
		int(duration.Hours()),
		int(duration.Minutes())%60,
		miner.Referred,
		miner.Active,
		address)

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
	// 		message = "To withdraw your funds, please open account on app.anotedao.com wallet and connect it to the bot!"
	// 	} else {
	// 		withdraw(miner.TelegramId)
	// 		message = "Your funds have been sent to your address. 🚀"
	// 	}
	// } else {
	// 	message = "You don't have enough funds to withdraw. The amount has to be bigger than 0.001 anotes. Please try later!"
	// }

	message := "This command has a bug and it has been temporarily disabled. Please try later!"
	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func alphaCommand(c telebot.Context) error {
	var err error
	message := "Exchange has been done successfully. Alpha version of Anote has been added to your beta balance in 1:10 ratio. 🚀"

	miner := getMiner(c.Message().Sender.ID)
	height := getHeight()

	if height <= 314000 {
		if strings.HasPrefix(miner.Address, "3A") {
			alp := &Alpha{}
			db.First(alp, &Alpha{Address: miner.Address})

			if alp.ID == 0 {
				ab := getAlphaBalance(miner.Address)
				if ab > 0 {
					log.Printf("Alpha balance: %s %d", miner.Address, ab)
					logTelegram(fmt.Sprintf("Alpha balance: %s %.8f", miner.Address, float64(ab)/10/float64(MULTI8)))
					alp.Address = miner.Address
					err := db.Save(alp).Error
					if err == nil {
						sendAsset(ab/10, "", miner.Address)
					}
				} else {
					alp.Address = miner.Address
					db.Save(alp)
					message = fmt.Sprintf("This address contains 0 anotes in alpha blockchain:\n\n%s", miner.Address)
				}
			} else {
				message = fmt.Sprintf("Alpha anotes from this address have already been exchanged:\n\n%s", miner.Address)
			}
		} else {
			message = fmt.Sprintf("The address is not valid:/n/n%s", miner.Address)
		}
	} else {
		message = "Exchange period ended with block 314000. You can check current block here:\n\nexplorer.anotedao.com"
	}

	_, err = bot2.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func checkCommand(c telebot.Context) error {
	var err error
	message := ""
	height := getHeight()

	if c.Message().IsReply() {
		miner := getMiner(c.Message().ReplyTo.Sender.ID)
		if miner.ID > 0 && miner.MiningHeight > 0 {
			diff := height - uint64(miner.MiningHeight)
			if diff <= 1410 {
				message = "This user is currently mining. 🚀"
			} else {
				message = fmt.Sprintf("This user is not mining currently, but has mined %d blocks ago.", diff)
			}
		} else {
			message = "This user never mined."
		}
		log.Println(prettyPrint(miner))
	} else {
		message = "To check the user, reply with /check to his message."
	}

	_, err = bot2.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func checkUserCommand(c telebot.Context) error {
	var err error
	m := c.Message()

	if m.IsForwarded() && m.Private() {
		tid := m.OriginalSender.ID
		m := getMiner(tid)

		message := fmt.Sprintf(`<b><u>Anote User Info</u></b>
		
Address: %s`,

			m.Address)

		_, err = bot2.Send(c.Chat(), message, telebot.NoPreview)
	}

	return err
}
