package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gopkg.in/telebot.v3"
)

type AnoteToday struct {
}

func (at *AnoteToday) sendAd(ad string) {
	var channelId int64
	if conf.Dev {
		channelId = TelDevAnoteToday
	} else {
		channelId = TelAnoteToday
	}
	r := &telebot.Chat{
		ID: channelId,
	}

	m, err := bot.Send(r, ad, telebot.NoPreview, telebot.Silent)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	num := int64(m.ID)
	dataTransaction2("%s__adnum", nil, &num, nil)
}

func (at *AnoteToday) start() {
	for {
		if at.isNewCycle() {
			code := at.generateNewCode()

			ad := at.getAd()

			at.sendAd(fmt.Sprintf(ad, code))
		}

		time.Sleep(time.Second * MonitorTick)
	}
}

func (at *AnoteToday) isNewCycle() bool {
	ks := &KeyValue{Key: "lastAdDay"}
	db.FirstOrCreate(ks, ks)
	today := time.Now().Day()

	if ks.ValueInt != uint64(today) && time.Now().Hour() == SendAdHour {
		ks.ValueInt = uint64(today)
		db.Save(ks)

		return true
	}

	return false
}

func (at *AnoteToday) generateNewCode() int {
	ks := &KeyValue{Key: "dailyCode"}
	db.FirstOrCreate(ks, ks)

	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 999

	code := rand.Intn(max-min+1) + min

	ks.ValueInt = uint64(code)
	db.Save(ks)

	return code
}

func (at *AnoteToday) getAd() string {
	ad := defaultAd5

	// cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// addr, err := proto.NewAddressFromString(TodayAddress)
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	// entries, _, err := cl.Addresses.AddressesData(ctx, addr)
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	// var winner *string
	// var amountWinner int64

	// for _, e := range entries {
	// 	if winner == nil {
	// 		address := parseItem(e.GetKey(), 0).(string)
	// 		amountWinner = e.ToProtobuf().GetIntValue()
	// 		winner = &address
	// 	} else {
	// 		amount := e.ToProtobuf().GetIntValue()
	// 		if amount > amountWinner {
	// 			amountWinner = amount
	// 			address := parseItem(e.GetKey(), 1).(string)
	// 			winner = &address
	// 		}
	// 	}
	// }

	// adData, err := getData(AdKey, winner)
	// if err != nil {
	// 	ad = defaultAd2
	// 	log.Println(err)
	// 	// logTelegram(err.Error())
	// } else {
	// 	adText := parseItem(adData.(string), 0)
	// 	adLink := parseItem(adData.(string), 1)
	// 	ad = adText.(string) + "\n\nRead <a href=\"" + adLink.(string) + "\">more</a>\n\n<a href=\"https://anotedao.com/advertise\"><strong><u>Advertise here!</u></strong></a>\n________________________\nDaily Mining Code: %d"

	// 	winnerKey := "%s__" + *winner
	// 	err := dataTransaction2(winnerKey, nil, nil, nil)
	// 	if err != nil {
	// 		log.Println(err)
	// 		logTelegram(err.Error())
	// 	}
	// }

	return ad
}

func initAnoteToday() {
	at := &AnoteToday{}
	go at.start()
}

var defaultAdBak = `<b><u>⭕️  ANOTE 2.0 IS NOW LIVE!</u></b>    🚀

We are proud to announce that Anote 2.0 is now available for mining.

We now have our own wallet (app.anotedao.com) which is used both as a wallet and a tool for mining. Stay tuned for more exciting news, information and tutorials!

You can find tutorial how to mine here: anotedao.com/mining

Join @AnoteDAO group for help and support!

________________________
@AnoteRobot Daily Mining Code: %d
`

var defaultAd = `Invite 3 or more users to Anote and your mining power will get multiplied by 10. You will get the referral link if you send /ref command to the miner.

You can find tutorial how to mine here: anotedao.com/mine

Join @AnoteDAO group for help and support!

<a href="https://anotedao.com/advertise"><strong><u>Advertise here!</u></strong></a>
________________________
@AnoteRobot Daily Mining Code: %d
`

var defaultAd1 = `Stake 10 AINTs or refer 3 new users and your mining power will get multiplied by 10.

You can find AINT tutorial here: anotedao.com/aint

<a href="https://anotedao.com/advertise"><strong><u>Advertise here!</u></strong></a>
________________________
@AnoteRobot Daily Mining Code: %d
`

var defaultAd2 = `Anote Staking Is Live 🚀

You can now stake your anotes and receive up to 70%% APY (Annual Percentage Yield). This means that if you stake 100 anotes, you will receive 70 more in a year. Payouts are done automatically every 10 minutes.

Read more about how to do it in the tutorial:

anotedao.com/staking

<a href="https://anotedao.com/advertise"><strong><u>Advertise here!</u></strong></a>
________________________
@AnoteRobot Daily Mining Code: %d
`

var defaultAd3 = `Mint AINT With Anotes 🚀

You can now use anotes to mint AINT which provides you with an opportunity for a long term investment because AINT behaves like a mining power for Anote.

Mint it by clicking "Mint AINT" tab in your wallet:

app.anotedao.com

<a href="https://anotedao.com/advertise"><strong><u>Advertise here!</u></strong></a>
________________________
@AnoteRobot Daily Mining Code: %d
`

var defaultAd4 = `Create Your Own Token 🚀

You can now create your own token on Anote chain and list it in Anote wallet. To create the token and list it, use our dev tool:

dev.anotedao.com

<a href="https://anotedao.com/advertise"><strong><u>Advertise here!</u></strong></a>
________________________
@AnoteRobot Daily Mining Code: %d
`

var defaultAd5 = `Mine More Anotes 🚀

<b><u>Additional 25%% for each referred user!</u></b>

Invite your friends to Anote and you will mine 25%% more aints every day for each referred user. You can find your referral link in your wallet (app.anotedao.com) or in menu if you open @AnoteRobot.

<a href="https://anotedao.com/advertise"><strong><u>Advertise here!</u></strong></a>
________________________
@AnoteRobot Daily Mining Code: %d
`
