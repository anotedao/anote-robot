package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	macaron "gopkg.in/macaron.v1"
	"gopkg.in/telebot.v3"
)

func viewTelegramLog(ctx *macaron.Context) {
	lr := NotificationResponse{}
	err := logTelegramService(ctx.Params("message"))
	lr.Success = err == nil
	ctx.JSON(200, lr)
}

func viewNotification(ctx *macaron.Context) {
	nr := &NotificationResponse{
		Success: true,
	}

	addr := ctx.Params("addr")
	h := ctx.Params("height")
	sh := ctx.Params("sheight")

	hi, err := strconv.Atoi(h)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		nr.Success = false
	}

	shi, err := strconv.Atoi(sh)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		nr.Success = false
	}

	minerData, err := getData(addr, nil)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		nr.Success = false
	}

	encId := parseItem(minerData.(string), 0)
	telId := DecryptMessage(encId.(string))

	idNum, err := strconv.Atoi(telId)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		nr.Success = false
	}

	rec := &telebot.Chat{
		ID: int64(idNum),
	}

	notification := "You have successfully started Anote mining cycle."

	if shi > 0 && hi-shi > 2880 {
		notification += "\n\n<u>Please notice that if you have continuity and mine on a daily basis, you receive a much bigger reward.</u>"
	}

	_, err = bot.Send(rec, notification)
	if err != nil {
		log.Println(err.Error() + " " + addr)
		logTelegram(err.Error() + " " + addr)
		nr.Success = false

		err := dataTransaction(addr, nil, nil, nil)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
	}

	ctx.JSON(200, nr)
}

func viewNotificationEnd(ctx *macaron.Context) {
	nr := &NotificationResponse{Success: true}
	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	adnum, err := getData2("%s__adnum", nil)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if err != nil {
		log.Println(err)
		nr.Success = false
	} else {
		message := fmt.Sprintf("Your mining cycle has ended.\n\nPlease run it again by getting the daily mining code in <a href=\"https://t.me/AnoteAds/%d\">AnoteAds</a> channel and sending it back here to reactivate mining cycle and withdraw already mined anotes. 🚀", adnum.(int64))
		rec := &telebot.Chat{
			ID: int64(tid),
		}
		bot.Send(rec, message, telebot.NoPreview)
	}

	ctx.JSON(200, nr)
}

func viewNotificationBattery(ctx *macaron.Context) {
	nr := &NotificationResponse{Success: true}
	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)

	if err != nil {
		log.Println(err)
		nr.Success = false
	} else {
		message := fmt.Sprint("<u><strong>Your AINT Miner health has dropped!</strong></u>❗️\n\nTo achieve 100% AINT Miner health and receive full amount of anotes, disable battery optimization on AINT Miner. You can learn how to do that here:\n\nanotedao.com/battery")
		rec := &telebot.Chat{
			ID: int64(tid),
		}
		bot.Send(rec, message, telebot.NoPreview)
	}

	ctx.JSON(200, nr)
}

func viewNotificationFirst(ctx *macaron.Context) {
	nr := &NotificationResponse{Success: true}
	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)
	if err != nil {
		log.Println(err)
		nr.Success = false
	} else {
		message := fmt.Sprint("<u><strong>You have successfully started Anote mining cycle!</strong></u>\n\nCheck your Anote balance by clicking the wallet button in the bottom left corner of your miner.\n\nYou will receive your first mining amount when you repeat the mining cycle after 24 hours. 🚀")
		rec := &telebot.Chat{
			ID: int64(tid),
		}
		bot.Send(rec, message)
	}

	ctx.JSON(200, nr)
}

type NotificationResponse struct {
	Success bool `json:"success"`
}

func inviteView(ctx *macaron.Context) {
	stats := getStats()

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctxb, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := proto.MustAddressFromString(MobileAddress)

	total, _, err := cl.Addresses.Balance(ctxb, addr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	basicAmount := float64(0)

	if stats.PayoutMiners > 0 {
		basicAmount = float64((total.Balance/(uint64(stats.PayoutMiners)+uint64(stats.ActiveReferred/4)))-Fee) / MULTI8
	} else {
		basicAmount = float64((total.Balance - Fee) / MULTI8)
	}

	nr := &NotificationResponse{Success: true}
	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)
	if err != nil {
		log.Println(err)
		nr.Success = false
	} else {
		da := pc.AnotePrice * basicAmount
		message := fmt.Sprintf("Please notice that by not mining Anote, you lose $%.2f daily.\n\nStart mining and earning today! 🚀", da)
		rec := &telebot.Chat{
			ID: int64(tid),
		}
		bot.Send(rec, message)
	}

	ctx.JSON(200, nr)
}
