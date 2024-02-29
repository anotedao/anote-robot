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
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

func initCommands() {
	bot2.Handle("/help", helpCommand2)
	bot2.Handle("/stats", statsCommand)
	bot2.Handle("/code", codeCommand)
	// bot2.Handle("/bo", batteryCommand)
	bot2.Handle("/alpha", alphaCommand)
	bot2.Handle("/bsc", addressBscCommand)
	bot2.Handle("/links", linksCommand)
	bot2.Handle("/seed", seedCommand)
	bot2.Handle("/check", checkCommand)
	bot2.Handle("/withdraw", withdrawCommandHelp)
	bot2.Handle(telebot.OnUserJoined, userJoined)
	bot2.Handle(telebot.OnText, checkUserCommand)
	// bot2.Handle(telebot.OnPhoto, addNewMessage)
	// bot2.Handle(telebot.OnMedia, addNewMessage)

	bot.Handle("/start", startCommand)
	bot.Handle("/miner", myStatsCommand)
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
			response = "You have successfully connected your Anote wallet. 🚀"
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

	s := fmt.Sprintf(`⭕️ <u><b>Anote Basic Stats</b></u>
	
	<b>Active Miners:</b> %d
	<b>Holders:</b> %d
	<b>Price:</b> $%.2f
	
	<b>Telegram Amount:</b> %.8f
	<b>sAINT Amount:</b> %.8f
	<b>Node Amount:</b> %.8f
	
	<b>Mined:</b> %s AINT
	<b>Community:</b> %s AINT
	<b>In Circulation:</b> %s AINT`,
		cch.StatsCache.ActiveMiners,
		cch.StatsCache.Holders,
		cch.StatsCache.Price,
		cch.StatsCache.AmountTlg,
		cch.StatsCache.AmountMobile,
		cch.StatsCache.AmountNode,
		cch.StatsCache.Mined,
		cch.StatsCache.Community,
		cch.StatsCache.Circulation)

	if m.Private() {
		s += fmt.Sprintf(`
	
	<b>Referred Miners:</b> %d
	<b>Payout Miners:</b> %d
	<b>Inactive Miners:</b> %d`,
			cch.StatsCache.Active,
			cch.StatsCache.Payout,
			cch.StatsCache.Inactive)
	}

	bot2.Send(m.Chat, s)

	return nil
}

func userJoined(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	msg := fmt.Sprintf("Hello, %s! Welcome to Anote community! 🚀\n\nHere are some resources to get you started:\n\nAnote Wallet: app.anotedao.com\nBlockchain Explorer: explorer.anotedao.com\nWebsite: anotedao.com\nMining Tutorial: anotedao.com/mine\nRun a Node: anotedao.com/node\n\n<u>Other Anote Communities:</u>\n\n@AnoteBalkan\n@AnoteAfrica\n@AnoteChina", m.Sender.FirstName)

	m, err := bot2.Send(m.Chat, msg, telebot.NoPreview)

	go func(m *telebot.Message) {
		time.Sleep(time.Second * 120)
		bot2.Delete(m)
	}(m)

	return err
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

	help := fmt.Sprintf("Daily mining code is at the bottom of the last announcement in <a href=\"https://t.me/AnoteAds/%d\"><strong><u>AnoteAds</u></strong></a> channel.", adnum.(int64))

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

func seedCommand(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	help := "<strong><u>Remember to save your seed!</u></strong>\n\nSeed is 15 words that you need to save somewhere on your phone, computer or write them on paper.\n\nYou can find your seed in the wallet settings when you click on the gear icon in the upper right corner of your wallet (app.anotedao.com)."

	_, err := bot2.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func withdrawCommandHelp(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	help := "*Can Anote be swapped to other cryptocurrencies?*\n\nYes, absolutely, Anote can be swapped to any other cryptocurrency by using PancakeSwap, Trust or MetaMask wallet and BSC chain. This is tutorial for using BSC gateway:\n\nanotedao.com/gateway\n\nThis is Anote token ID / address in BSC chain:\n\n`0xe7f0f1585bdbd06b18dbb87099b87bd79bbd315b`\n\nThis information is for crypto experts. If this is not enough for you, please wait for active development to be finished and more tutorials to be made. Until then, mine AINT, mint ANOTE with it and stake it."

	_, err := bot2.Reply(m, help, telebot.ModeMarkdown, telebot.NoPreview)
	// _, err := bot2.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func addressBscCommand(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	help := "Address of Anote contract in BSC chain:\n\n`0xe7f0f1585bdbd06b18dbb87099b87bd79bbd315b`"

	// o := telebot.SendOptions{
	// 	ParseMode:             telebot.ModeMarkdownV2,
	// 	DisableWebPagePreview: true,
	// }
	_, err := bot2.Send(m.Chat, help, telebot.ModeMarkdownV2)

	return err
}

func linksCommand(c telebot.Context) error {
	saveUser(c)
	m := c.Message()

	help := `⭕️ <b><u>Important Anote Links</u></b>

Wallet - app.anotedao.com
Website - anotedao.com
BSC Gateway - anotedao.com/gateway`

	_, err := bot2.Send(m.Chat, help, telebot.NoPreview)

	return err
}

func mineCommand(c telebot.Context) error {
	var err error
	m := c.Message()
	log.Println(prettyPrint(m))
	log.Println(m.IsForwarded())
	if c.Message().Private() {
		if c.Message().IsForwarded() {
			return checkUserCommand(c)
		} else {
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

	// <b>Mined Telegram:</b> %.8f ANOTE
	// <b>Mined Mobile:</b> N/A
	// float64(miner.MinedTelegram)/float64(MULTI8),

	message := fmt.Sprintf(`⭕️ <b><u>Your Anote Stats</u></b>

	<b>Address Balance:</b> %.8f AINT

	<b>Cycle Blocks Left:</b> %d
	<b>Cycle Time Left:</b> %02d:%02d

	<b>Referred Total:</b> %d
	<b>Referred Active:</b> %d

	<b><u>Address</u></b>

	%s`,
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

	message = fmt.Sprintf("https://anotedao.com/mine?r=%d", miner.ID)
	_, err = bot.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func alphaCommand(c telebot.Context) error {
	var err error
	message := "Exchange has been done successfully. Alpha version of Anote has been added to your beta balance in 1:10 ratio. 🚀"

	miner := getMiner(c.Message().Sender.ID)
	height := getHeight()

	if height <= 560000 {
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
		message = "Exchange period ended with block 560000. You can check current block here:\n\nexplorer.anotedao.com"
	}

	_, err = bot2.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func checkCommand(c telebot.Context) error {
	var err error
	var miner *MinerResponse
	message := ""
	height := getHeight()

	ks := &KeyValue{Key: "dailyCode"}
	db.FirstOrCreate(ks, ks)

	code := ks.ValueInt

	if c.Message().IsReply() {
		miner = getMiner(c.Message().ReplyTo.Sender.ID)
		if miner.ID > 0 && miner.MiningHeight > 0 {
			diff := height - uint64(miner.MiningHeight)
			if diff <= 1410 {
				message = "This user is currently mining. 🚀"
			} else {
				message = fmt.Sprintf("This user is not mining currently, but has mined %d blocks ago.\n\nTo continue mining, send %d as a message to @AnoteRobot.", diff, code)
			}
		} else {
			message = fmt.Sprintf("This user has never mined.\n\nTo start mining, send %d as a message to @AnoteRobot.", code)
		}
		log.Println(prettyPrint(miner))
	} else {
		miner = getMiner(c.Message().Sender.ID)
		if miner.ID > 0 && miner.MiningHeight > 0 {
			diff := height - uint64(miner.MiningHeight)
			if diff <= 1410 {
				message = "You are currently mining. 🚀"
			} else {
				message = fmt.Sprintf("You are not mining currently, but you have mined %d blocks ago.\n\nTo continue mining, send %d as a message to @AnoteRobot.", diff, code)
			}
		} else {
			message = fmt.Sprintf("You haven't mined so far.\n\nTo start mining, send %d as a message to @AnoteRobot.", code)
		}
	}

	log.Println(prettyPrint(miner))

	_, err = bot2.Send(c.Chat(), message, telebot.NoPreview)

	return err
}

func checkUserCommand(c telebot.Context) error {
	var err error
	m := c.Message()

	log.Println(prettyPrint(m))

	if m.IsForwarded() && m.Private() {
		if m.OriginalSender != nil {
			tid := m.OriginalSender.ID
			m := getMiner(tid)

			message := fmt.Sprintf(`<b><u>Anote User Info</u></b>
			
	Address: %s`,

				m.Address)

			_, err = bot2.Send(c.Chat(), message, telebot.NoPreview)
		}
	} else if !m.Private() {
		if m.Chat.ID == TelAnon {
			newMessage(m)
		}
		txt := m.Text
		if len(txt) == 3 {
			code, err := strconv.Atoi(txt)
			if err == nil && code < 1000 {
				bot2.Reply(m, "Please send the code to @AnoteRobot!", telebot.NoPreview)
			}
		}

		group, err := bot2.ChatByID(m.Chat.ID)
		if err != nil {
			log.Println(err)
		}

		cm, err := bot2.ChatMemberOf(group, m.Sender)
		if err != nil {
			log.Println(err)
		}

		log.Println(prettyPrint(cm))

		if (cm.Role != telebot.Administrator && cm.Role != telebot.Creator) &&
			m.Chat.ID != TelGroup &&
			(strings.Contains(strings.ToLower(txt), strings.ToLower("withdraw")) ||
				strings.Contains(strings.ToLower(txt), strings.ToLower("swap")) ||
				strings.Contains(strings.ToLower(txt), strings.ToLower("exchange")) ||
				strings.Contains(strings.ToLower(txt), strings.ToLower("buy")) ||
				strings.Contains(strings.ToLower(txt), strings.ToLower("cash")) ||
				strings.Contains(strings.ToLower(txt), strings.ToLower("dump")) ||
				strings.Contains(strings.ToLower(txt), strings.ToLower("sell"))) {
			withdrawCommandHelp(c)
		}
	}

	return err
}

func addNewMessage(c telebot.Context) error {
	m := c.Message()

	log.Println(prettyPrint(m))

	if m.Chat.ID == TelAnoteNews {
		newMessage(m)
	}

	return nil
}
