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
	m.Get("/notification-weekly/:telegramid", viewNotificationWeekly)
	m.Get("/notification-bo/:telegramid", viewNotificationBattery)
	m.Get("/notification-first/:telegramid", viewNotificationFirst)
	m.Get("/notification-tg/:telegramid/:message", viewNotificationTelegram)
	m.Get("/is-follower/:telegramid", viewIsFollower)
	m.Get("/alpha-sent/:addr", viewAlphaSent)

	go m.Run("0.0.0.0", Port)

	return m
}
