package main

import (
	"log"
	"strconv"

	macaron "gopkg.in/macaron.v1"
	"gopkg.in/telebot.v3"
)

func viewNotification(ctx *macaron.Context) {
	nr := &NotificationResponse{
		Success: true,
	}

	addr := ctx.Params("addr")
	ta := TelegramAddress

	encId, err := getData(addr, &ta)
	if err != nil {
		log.Println(err)
		nr.Success = false
	}
	telId := DecryptMessage(encId.(string))

	idNum, err := strconv.Atoi(telId)
	if err != nil {
		log.Println(err)
		nr.Success = false
	}

	rec := &telebot.Chat{
		ID: int64(idNum),
	}

	notification := "You have successfully started Anote mining cycle."

	_, err = bot.Send(rec, notification)
	if err != nil {
		log.Println(err)
		nr.Success = false
	}

	ctx.JSON(200, nr)
}

type NotificationResponse struct {
	Success bool `json:"success"`
}
