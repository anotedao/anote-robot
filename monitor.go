package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/telebot.v3"
)

type Monitor struct {
	Miners *MinersResponse
	Height uint64
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
		(int(m.Height)-miner.MiningHeight) > 1440 &&
		(int(m.Height)-miner.MiningHeight) < 2880 &&
		time.Since(dbminer.LastNotification) > time.Hour*24 {

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

func initMonitor() *Monitor {
	m := &Monitor{
		Miners: &MinersResponse{},
	}
	go m.start()
	return m
}

type MinersResponse []*MinerResponse

type MinerResponse struct {
	Address          string    `json:"Address"`
	LastNotification time.Time `json:"LastNotification"`
	TelegramID       int       `json:"TelegramId"`
	MiningHeight     int       `json:"MiningHeight"`
}
