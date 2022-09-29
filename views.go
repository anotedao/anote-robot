package main

import (
	"log"

	macaron "gopkg.in/macaron.v1"
)

func viewNotification(ctx *macaron.Context) {
	nr := &NotificationResponse{
		Success: true,
		Error:   0,
	}

	log.Println(ctx.Params("addr"))

	ctx.JSON(200, nr)
}

type NotificationResponse struct {
	Success bool `json:"success"`
	Error   int  `json:"error"`
}
