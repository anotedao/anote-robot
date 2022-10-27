package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Prices struct {
	BTC float64 `json:"BTC"`
	ETH float64 `json:"ETH"`
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
}

type PriceClient struct {
	Prices     *Prices
	Loaded     bool
	AnotePrice float64
}

func (pc *PriceClient) doRequest() (*Prices, error) {
	p := &Prices{}
	cl := http.Client{}

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, PricesURL, nil)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)

	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			log.Println(string(body))
			err := errors.New(res.Status)
			return nil, err
		}
		json.Unmarshal(body, p)
	} else {
		return nil, err
	}

	pc.Prices = p

	pc.doRequestOrderbook()

	return p, nil
}

func (pc *PriceClient) doRequestOrderbook() {
	or := &OrderbookResponse{}
	cl := http.Client{}

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, OrderbookURL, nil)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return
	}

	res, err := cl.Do(req)

	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
			return
		}
		if res.StatusCode != 200 {
			err := errors.New(res.Status)
			log.Println(err)
			logTelegram(err.Error())
			return
		}
		json.Unmarshal(body, or)
	} else {
		log.Println(err)
		logTelegram(err.Error())
		return
	}

	pc.AnotePrice = float64(or.Asks[0].Price) / MULTI8 * pc.Prices.USD
}

func (pc *PriceClient) start() {
	go func() {
		for {
			if p, err := pc.doRequest(); err != nil {
				log.Println(err)
				logTelegram(err.Error())
			} else {
				pc.Prices = p
			}

			// if conf.Debug {
			// 	log.Printf("%#v\n", pc.Prices)
			// }

			pc.Loaded = true

			time.Sleep(time.Minute * 5)
		}
	}()
}

func initPriceClient() *PriceClient {
	pc := &PriceClient{Loaded: false}
	pc.start()
	return pc
}

type OrderbookResponse struct {
	Timestamp int64 `json:"timestamp"`
	Pair      struct {
		AmountAsset string `json:"amountAsset"`
		PriceAsset  string `json:"priceAsset"`
	} `json:"pair"`
	Bids []struct {
		Amount int64 `json:"amount"`
		Price  int   `json:"price"`
	} `json:"bids"`
	Asks []struct {
		Amount int64 `json:"amount"`
		Price  int   `json:"price"`
	} `json:"asks"`
}