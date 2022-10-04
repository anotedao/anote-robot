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
	ctx.JSON(200, lr)
}

func viewNotification(ctx *macaron.Context) {
	nr := &NotificationResponse{
		Success: true,
	}

	addr := ctx.Params("addr")

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

	_, err = bot.Send(rec, notification)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		nr.Success = false
	}

	ctx.JSON(200, nr)
}

type NotificationResponse struct {
	Success bool `json:"success"`
}
