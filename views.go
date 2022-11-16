package main

import (
	"log"
	"strconv"

	macaron "gopkg.in/macaron.v1"
	"gopkg.in/telebot.v3"
)

func viewTelegramLog(ctx *macaron.Context) {
	lr := NotificationResponse{}
	err := logTelegramService(ctx.Params("message"))
	lr.Success = err == nil
	log.Println(err)
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

type NotificationResponse struct {
	Success bool `json:"success"`
}
