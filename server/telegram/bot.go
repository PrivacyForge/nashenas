package telegram

import (
	"log"

	"github.com/PrivacyForge/nashenas/configs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func InitBot() tgbotapi.UpdatesChannel {

	var err error
	Bot, err = tgbotapi.NewBotAPI(configs.BotToken)

	if err != nil {
		panic("telegram bot connection failed.")
	}

	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	return updates
}
