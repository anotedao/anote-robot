package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

type Monitor struct {
	Miners proto.DataEntries
}

func (m *Monitor) loadMiners() {
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	addr, err := proto.NewAddressFromPublicKey(55, key)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	m.Miners, _, err = cl.Addresses.AddressesData(ctx, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
}

func (m *Monitor) sendNotifications() {
	for _, miner := range m.Miners {
		if m.isSending(miner) {
			m.sendNotification(miner)
			log.Println(miner.GetKey())
		}
	}
}

func (m *Monitor) isSending(miner proto.DataEntry) bool {
	key := miner.GetKey()
	height := getHeight()
	mobile := MobileAddress

	minerHeight, _ := getData(key, &mobile)
	// if err != nil {
	// 	log.Println(err)
	// logTelegram(err.Error())
	// }

	dbminer := &Miner{}
	db.FirstOrCreate(dbminer, Miner{Address: key})

	// log.Printf("%s %d %d", key, minerHeight.(int64), int64(height)-minerHeight.(int64))

	// log.Println(prettyPrint(dbminer))

	if minerHeight != nil &&
		(int64(height)-minerHeight.(int64)) > 1440 &&
		(int64(height)-minerHeight.(int64)) < 2880 &&
		time.Since(dbminer.LastNotification) > time.Hour*24 {

		dbminer.LastNotification = time.Now()
		db.Save(dbminer)

		return true
	}

	return false
}

func (m *Monitor) sendNotification(miner proto.DataEntry) {
	notification := fmt.Sprint("Your mining period has ended. Please run it again to reactivate and withdraw already mined anotes. 🚀\n\nYou can find daily mining code in @AnoteToday channel.")

	encId := miner.ToProtobuf().GetStringValue()
	telId := DecryptMessage(encId)

	idNum, err := strconv.Atoi(telId)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	rec := &telebot.Chat{
		ID: int64(idNum),
	}

	_, err = bot.Send(rec, notification)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
}

func (m *Monitor) minerExists(telId string) bool {
	for _, mnr := range m.Miners {
		mnrTelId := DecryptMessage(mnr.ToProtobuf().GetStringValue())
		if mnrTelId == telId {
			return true
		}
	}

	return false
}

func (m *Monitor) start() {
	m.loadMiners()

	go func() {
		for {
			m.loadMiners()
			time.Sleep(time.Hour)
		}
	}()

	for {
		m.sendNotifications()

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	return m
}
