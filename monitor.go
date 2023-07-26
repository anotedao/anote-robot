package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Monitor struct {
	Height             uint64
	BeneficiaryBalance uint64
}

func (m *Monitor) monitorAintBuys() {
	for {
		cl, err := client.NewClient(client.Options{BaseUrl: WavesNodeURL, Client: &http.Client{}})
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)

		addr := proto.MustAddressFromString(conf.Beneficiary)

		total, _, err := cl.Addresses.Balance(ctx, addr)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		} else {
			if total.Balance > m.BeneficiaryBalance {
				nb := float64(total.Balance-m.BeneficiaryBalance) / MULTI8
				notificationTelegram(fmt.Sprintf("<u><strong>New AINT minted:</strong></u> 🚀\n\n%.8f WAVES", nb))
				notificationTelegramTeam(fmt.Sprintf("<u><strong>New AINT minted:</strong></u> 🚀\n\n%.8f WAVES", nb))
			}

			m.BeneficiaryBalance = total.Balance
		}

		cancel()

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
