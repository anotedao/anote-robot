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
				usd := nb * pc.AnotePrice
				// notificationTelegram(fmt.Sprintf("<u><strong>New AINT Minted!</strong></u> 🚀\n\nPaid:\n%.8f ANOTE ($%.2f)\nMinted:\n%.8f AINT", nb, usd, naints))
				// notificationTelegramTeam(fmt.Sprintf("<u><strong>New AINT Minted!</strong></u> 🚀\n\nPaid:\n%.8f ANOTE ($%.2f)\nMinted:\n%.8f AINT", nb, usd, naints))
				notificationTelegramGroup(fmt.Sprintf("<u><strong>New AINT Minted!</strong></u> 🚀\n\nPaid:\n%.8f ANOTE ($%.2f)\nMinted:\n%.8f AINT", nb, usd, naints))
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
				notificationTelegramTeamPin(fmt.Sprintf("<u><strong>AINT Price Increased!</strong></u> 🚀\n\nNew Price:\n$%.2f", apf))
				notificationTelegramGroupPin(fmt.Sprintf("<u><strong>AINT Price Increased!</strong></u> 🚀\n\nNew Price:\n$%.2f", apf))
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
		nnodes := uint64(0)
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
				nnodes = m.NodeBalance - nodes.Balance
			}
			m.NodeBalance = nodes.Balance
		}

		if nnodes > 0 && count > 0 {
			notificationTelegram(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token", nnodes))
			notificationTelegramTeam(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token", nnodes))
			notificationTelegramGroup(fmt.Sprintf("<u><strong>New NODE Minted!</strong></u> 🚀\n\n%d NODE\n\nAbout NODE Token:\nanotedao.com/node-token", nnodes))
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

func initMonitor() *Monitor {
	m := &Monitor{}
	// go m.start()
	go m.monitorAintBuys()
	go m.monitorNodeMints()
	go m.monitorDiskSpace()
	return m
}
