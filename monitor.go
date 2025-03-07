package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

type Monitor struct {
	Height             uint64
	BeneficiaryBalance uint64
	NodeBalance        uint64
	AintBalance        uint64
	AintPrice          float64
	NodePrice          float64
}

func (m *Monitor) monitorAintBuys() {
	count := 0

	asset, err := crypto.NewDigestFromBase58(AintAnoteId)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	for {
		naints := float64(0)
		cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)

		addr := proto.MustAddressFromString(conf.Beneficiary)

		aints, _, err := cl.Assets.BalanceByAddressAndAsset(ctx, addr, asset)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		} else {
			if m.AintBalance != aints.Balance {
				naints = float64(m.AintBalance-aints.Balance) / MULTI8
			}
			m.AintBalance = aints.Balance
		}

		total, _, err := cl.Addresses.Balance(ctx, addr)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		} else {
			if total.Balance > m.BeneficiaryBalance && count > 0 {
				nb := float64(total.Balance-m.BeneficiaryBalance) / MULTI8
				// usd := nb * pc.AnotePrice
				// notificationTelegram(fmt.Sprintf("<u><strong>New AINT Minted!</strong></u> 🚀\n\nPaid:\n%.8f ANOTE ($%.2f)\nMinted:\n%.8f AINT", nb, usd, naints))
				// notificationTelegramTeam(fmt.Sprintf("<u><strong>New ANOTE Minted!</strong></u> 🚀\n\nPaid:\n%.8f AINT\nMinted:\n%.8f ANOTE", nb, naints))
				notificationTelegramGroup(fmt.Sprintf("<u><strong>New ANOTE Minted!</strong></u> 🚀\n\nPaid:\n%.8f AINT\nMinted:\n%.8f ANOTE", nb, naints))
				// notificationTelegramGroupBalkan(fmt.Sprintf("<u><strong>New ANOTE Minted!</strong></u> 🚀\n\nPaid:\n%.8f AINT\nMinted:\n%.8f ANOTE", nb, naints))
			}

			m.BeneficiaryBalance = total.Balance
		}

		if count > 0 {
			ap, err := getData2("%s__price", nil)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			apf := float64(ap.(int64)) / MULTI8

			if apf > m.AintPrice {
				notificationTelegramTeamPin(fmt.Sprintf("<u><strong>ANOTE Price Increased!</strong></u> 🚀\n\nNew Price:\n%.2f AINT", apf))
				notificationTelegramGroupPin(fmt.Sprintf("<u><strong>ANOTE Price Increased!</strong></u> 🚀\n\nNew Price:\n%.2f AINT", apf))
				// notificationTelegramGroupBalkanPin(fmt.Sprintf("<u><strong>ANOTE Price Increased!</strong></u> 🚀\n\nNew Price:\n%.2f AINT", apf))
			}
		}

		ap, err := getData2("%s__price", nil)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
		m.AintPrice = float64(ap.(int64)) / MULTI8

		cancel()

		count++

		time.Sleep(time.Second * 60)
	}
}

func (m *Monitor) monitorNodeMints() {
	count := 0

	asset, err := crypto.NewDigestFromBase58(NodeAnoteId)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	for {
		nnodes := int64(0)
		cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)

		addr := proto.MustAddressFromString(GATEWAY_ADDR)

		nodes, _, err := cl.Assets.BalanceByAddressAndAsset(ctx, addr, asset)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		} else {
			if m.NodeBalance != nodes.Balance {
				nnodes = int64(m.NodeBalance) - int64(nodes.Balance)
			}
			m.NodeBalance = nodes.Balance
		}

		if nnodes > 0 && count > 0 {
			ga := GATEWAY_ADDR
			nt, err := getData("%s__nodeTier", &ga)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			notificationTelegram(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token\n\n<strong><u>%d NODE left at the price of %.2f BNB.</u></strong>", nnodes, nt.(int64), m.NodePrice))
			notificationTelegramTeam(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token\n\n<strong><u>%d NODE left at the price of %.2f BNB.</u></strong>", nnodes, nt.(int64), m.NodePrice))
			notificationTelegramGroup(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token\n\n<strong><u>%d NODE left at the price of %.2f BNB.</u></strong>", nnodes, nt.(int64), m.NodePrice))
			// notificationTelegramGroupBalkan(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token\n\n<strong><u>%d NODE left at the price of %.2f BNB.</u></strong>", nnodes, nt.(int64), m.NodePrice))
		}

		if count > 0 {
			ga := GATEWAY_ADDR
			np, err := getData("%s__nodePrice", &ga)
			if err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			npf := float64(np.(int64)) / 100

			if npf > m.NodePrice {
				notificationTelegramTeamPin(fmt.Sprintf("<u><strong>NODE Price Increased!</strong></u> 🚀\n\nNew Price:\n$%.2f BNB", npf))
				notificationTelegramGroupPin(fmt.Sprintf("<u><strong>NODE Price Increased!</strong></u> 🚀\n\nNew Price:\n$%.2f BNB", npf))
				// notificationTelegramGroupBalkanPin(fmt.Sprintf("<u><strong>NODE Price Increased!</strong></u> 🚀\n\nNew Price:\n$%.2f BNB", npf))
			}
		}

		ga := GATEWAY_ADDR
		np, err := getData("%s__nodePrice", &ga)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
		m.NodePrice = float64(np.(int64)) / 100

		cancel()

		count++

		time.Sleep(time.Second * 60)
	}
}

func (m *Monitor) monitorDiskSpace() {
	for {
		time.Sleep(time.Second * 30)
	}
}

func (m *Monitor) forwardCompetition() {
	for {
		group, err := bot2.ChatByID(TelAnon)
		if err != nil {
			log.Println(err)
		}

		ch, err := bot2.ChatByID(TelAnoteNews)
		if err != nil {
			log.Println(err)
		}

		msg := &telebot.Message{}
		msg.ID = 17
		msg.Chat = ch

		if m1 != nil &&
			m2 != nil &&
			m3 != nil &&
			msg.ID != m1.ID &&
			msg.ID != m2.ID &&
			msg.ID != m3.ID {

			bot2.Forward(group, msg, telebot.NoPreview)
			newMessage(msg)
		}
		time.Sleep(time.Minute * 5)
	}
}

func (m *Monitor) monitorNodes() {
	for {
		cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		addr := proto.MustAddressFromString(COMMUNITY_ADDR)

		data, resp, err := cl.Addresses.AddressesData(ctx, addr)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
		resp.Body.Close()

		for _, de := range data {
			if !isNodeActive(de.GetKey()) {
				sentNodeNotification(de.GetKey(), de.ToProtobuf().GetStringValue())
			}
		}

		time.Sleep(time.Second * 30)

		cancel()
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	// go m.start()
	go m.monitorAintBuys()
	go m.monitorNodeMints()
	go m.monitorDiskSpace()
	go m.forwardCompetition()
	go m.monitorNodes()
	return m
}

type StoredMessage struct {
	MessageID int
	ChatID    int64
}

func (sm StoredMessage) MessageSig() (int, int64) {
	return sm.MessageID, sm.ChatID
}
