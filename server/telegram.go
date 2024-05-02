package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func InitBot() (tgbotapi.UpdatesChannel, error) {
	BOT_TOKEN := os.Getenv("BOT_TOKEN")

	var err error
	bot, err = tgbotapi.NewBotAPI(BOT_TOKEN)

	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return updates, nil
}
