package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Cache struct {
	StatsCache *StatsCache
}

func (c *Cache) loadStatsCache() {
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

	abr4, err := anc.AddressesBalance(COMMUNITY_ADDR2)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	abr5, err := anc.AddressesBalance(AINT_MINT_ADDR)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	abr6, err := anc.AddressesBalance(ANOTE_STAKE_ADDR)
	if err != nil {
		log.Println(err.Error())
		logTelegram(err.Error())
	}

	balance := (abr.Balance / int(SATINBTC)) + (abr2.Balance / int(SATINBTC)) + (abr3.Balance / int(SATINBTC)) + (abr4.Balance / int(SATINBTC)) + (abr5.Balance / int(SATINBTC)) + (abr6.Balance / int(SATINBTC))
	circulation := mined - int64(balance)

	stats := getStats()

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	addr := proto.MustAddressFromString(MobileAddress)

	total, _, err := cl.Addresses.Balance(ctx, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	// addrT := proto.MustAddressFromString(TelegramAddress)

	// totalT, _, err := cl.Addresses.Balance(ctx, addr)
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	basicAmount := float64(0)
	basicAmountT := float64(0)

	if stats.ActiveUnits > 0 {
		basicAmount = float64((float64(total.Balance) / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4)))) / MULTI8
		basicAmountT = float64((43.2 / float64(uint64(stats.ActiveUnits)+uint64(stats.ActiveReferred/4))))
	} else {
		basicAmount = float64((float64(total.Balance)) / MULTI8)
		basicAmountT = 43.2
	}

	c.StatsCache.ActiveMiners = stats.ActiveMiners
	// c.StatsCache.Holders = stats.Holders
	c.StatsCache.Price = pc.AnotePrice
	c.StatsCache.AmountTlg = basicAmountT
	c.StatsCache.AmountMobile = basicAmount
	c.StatsCache.Mined = humanize.Comma(mined)
	c.StatsCache.Community = humanize.Comma(int64(balance))
	c.StatsCache.Circulation = humanize.Comma(circulation)
	c.StatsCache.Active = stats.ActiveReferred
	c.StatsCache.Payout = stats.PayoutMiners
	c.StatsCache.Inactive = stats.InactiveMiners
	c.StatsCache.Referred = stats.ActiveReferred

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get("http://localhost:5005/distribution")
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var ds DistributionResponse
	if err := json.Unmarshal(body, &ds); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	c.StatsCache.Holders = len(ds)

	cancel()
}

func (c *Cache) start() {
	for {
		c.loadStatsCache()

		time.Sleep(time.Second * 10)
	}
}

func initCache() *Cache {
	c := &Cache{}
	c.StatsCache = &StatsCache{}
	go c.start()

	return c
}

type StatsCache struct {
	ActiveMiners int
	Holders      int
	Price        float64
	AmountTlg    float64
	AmountMobile float64
	Mined        string
	Community    string
	Circulation  string
	Active       int
	Payout       int
	Inactive     int
	Referred     int
}
