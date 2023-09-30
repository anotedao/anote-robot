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

	notification := "<b><u>You have successfully started Anote mining cycle!</u></b> 🚀"

	if shi > 0 && hi-shi > 2880 {
		notification += "\n\n<u>Please notice that if you have continuity and mine on a daily basis, you receive a much bigger reward.</u>"
	}

	notification += `<b><u>Commands:</u></b>

	/miner - Your miner stats
	/ref - Your Anote referral link
	/withdraw - Withdraw your mined anotes`

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
		message := fmt.Sprintf("Your mining cycle has ended.\n\nPlease run it again by getting the daily mining code in <a href=\"https://t.me/AnoteAds/%d\">AnoteAds</a> channel and sending it back here to reactivate the mining cycle. 🚀", adnum.(int64))
		rec := &telebot.Chat{
			ID: int64(tid),
		}
		bot.Send(rec, message, telebot.NoPreview)
	}

	ctx.JSON(200, nr)
}

func viewNotificationWeekly(ctx *macaron.Context) {
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

	stats := cch.StatsCache

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

	if stats.Payout > 0 {
		basicAmount = float64((total.Balance/(uint64(stats.Payout)+uint64(stats.Referred/4)))-Fee) / MULTI8
	} else {
		basicAmount = float64((total.Balance - Fee) / MULTI8)
	}

	if err != nil {
		log.Println(err)
		nr.Success = false
	} else {
		message := fmt.Sprintf("<b><u>Telegram mining is back!</u></b> 🚀\n\nStart mining by getting the daily mining code in <a href=\"https://t.me/AnoteAds/%d\">AnoteAds</a> channel and sending it back here to reactivate the mining cycle.\n\nJoin @AnoteDAO group for help and support!", adnum.(int64))
		da := pc.AnotePrice * basicAmount
		message += fmt.Sprintf("\n\nPlease notice that by not mining Anote, you lose $%.2f daily.\n\nStart mining and earning today! 🚀", da)
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
		message := "<u><strong>You have successfully started Anote mining cycle!</strong></u> 🚀\n\nCheck your Anote balance with /miner command!"
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
	nr := &NotificationResponse{Success: false}
	stats := cch.StatsCache
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

	if stats.Payout > 0 {
		basicAmount = float64((total.Balance/(uint64(stats.Payout)+uint64(stats.Referred/4)))-Fee) / MULTI8
	} else {
		basicAmount = float64((total.Balance - Fee) / MULTI8)
	}

	tids := ctx.Params("telegramid")
	tid, err := strconv.Atoi(tids)
	if err != nil {
		log.Println(err)
	} else {
		nr.Success = true
		da := pc.AnotePrice * basicAmount
		message := fmt.Sprintf("Please notice that by not mining Anote, you lose $%.2f daily.\n\nStart mining and earning today! 🚀", da)
		rec := &telebot.Chat{
			ID: int64(tid),
		}
		bot.Send(rec, message)
	}

	ctx.JSON(200, nr)
}

func viewAlphaSent(ctx *macaron.Context) {
	alr := &AlphaSentResponse{Sent: false}
	addr := ctx.Params("addr")

	if len(addr) > 0 {
		alp := &Alpha{}
		db.First(alp, &Alpha{Address: addr})
		if alp.ID > 0 {
			alr.Sent = true
		}
	}

	ctx.JSON(200, alr)
}

type AlphaSentResponse struct {
	Sent bool `json:"sent"`
}
