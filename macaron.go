package main

import (
	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(macaron.Renderer())
	m.Use(cache.Cacher())

	m.Get("/notification/:addr/:height/:sheight", viewNotification)
	m.Get("/log/:message", viewTelegramLog)
	m.Get("/invite/:telegramid", inviteView)
	m.Get("/notification-end/:telegramid", viewNotificationEnd)
	m.Get("/notification-bo/:telegramid", viewNotificationBattery)

	go m.Run("127.0.0.1", Port)

	return m
}
