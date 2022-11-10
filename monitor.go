package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

type Monitor struct {
	Miners             *MinersResponse
	Height             uint64
	BeneficiaryBalance uint64
}

func (m *Monitor) loadMiners() {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5003/miners"))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, m.Miners); err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
	}
}

func (m *Monitor) sendNotifications() {
	for _, miner := range *m.Miners {
		if m.isSending(miner) {
			m.sendNotification(miner)
			log.Printf("Notification: %s", miner.Address)
		}
	}
}

func (m *Monitor) isSending(miner *MinerResponse) bool {
	dbminer := &Miner{}
	tx := db.FirstOrCreate(dbminer, Miner{Address: miner.Address})

	if dbminer.ID != 0 &&
		tx.Error == nil &&
		(int(m.Height)-miner.MiningHeight) > 1410 &&
		(int(m.Height)-miner.MiningHeight) < 2000 &&
		dbminer.LastNotification.Day() != time.Now().Day() {

		dbminer.LastNotification = time.Now()
		db.Save(dbminer)

		return true
	}

	return false
}

func (m *Monitor) sendNotification(miner *MinerResponse) {
	notification := fmt.Sprint("Your mining period has ended. Please run it again to reactivate and withdraw already mined anotes. 🚀\n\nYou can find daily mining code in @AnoteToday channel.")

	telId := miner.TelegramID

	rec := &telebot.Chat{
		ID: int64(telId),
	}

	_, err := bot.Send(rec, notification)
	if err != nil {
		log.Println(err.Error() + " " + miner.Address)
		logTelegram(err.Error() + " " + miner.Address)
	}
}

func (m *Monitor) minerExists(telId int64) bool {
	for _, mnr := range *m.Miners {
		if int64(mnr.TelegramID) == telId {
			return true
		}
	}

	return false
}

func (m *Monitor) start() {
	m.loadMiners()

	go func() {
		for {
			m.Height = getHeight()
			m.loadMiners()
			time.Sleep(time.Second * 30)
		}
	}()

	for {
		m.sendNotifications()

		time.Sleep(time.Second * MonitorTick)
	}
}

func (m *Monitor) monitorAintBuys() {
	for {
		cl, err := client.NewClient(client.Options{BaseUrl: WavesNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		addr := proto.MustAddressFromString(conf.Beneficiary)

		total, _, err := cl.Addresses.Balance(ctx, addr)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		if total.Balance > m.BeneficiaryBalance {
			nb := float64(total.Balance-m.BeneficiaryBalance) / MULTI8
			notificationTelegram(fmt.Sprintf("<u><strong>New AINT purchase: %.8f WAVES</strong></u> 🚀", nb))
		}

		m.BeneficiaryBalance = total.Balance

		time.Sleep(time.Second * 10)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{
		Miners: &MinersResponse{},
	}
	go m.start()
	go m.monitorAintBuys()
	return m
}

type MinersResponse []*MinerResponse

type MinerResponse struct {
	Address          string    `json:"address"`
	LastNotification time.Time `json:"last_notification"`
	TelegramID       int       `json:"telegram_id"`
	MiningHeight     int       `json:"mining_height"`
}
