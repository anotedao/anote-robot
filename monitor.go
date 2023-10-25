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
	AintBalance        uint64
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
				notificationTelegram(fmt.Sprintf("<u><strong>New AINT minted!</strong></u> 🚀\n\n$%.8f\n%.8f ANOTE\n%.8f AINT", usd, nb, naints))
				// notificationTelegramTeam(fmt.Sprintf("<u><strong>New AINT minted!</strong></u> 🚀\n\n%.8f ANOTE", nb))
			}

			m.BeneficiaryBalance = total.Balance
		}

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
	go m.monitorDiskSpace()
	return m
}
