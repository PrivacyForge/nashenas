package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"

	"github.com/PrivacyForge/nashenas/configs"
	"github.com/PrivacyForge/nashenas/database"
	"github.com/PrivacyForge/nashenas/routes"
	"github.com/PrivacyForge/nashenas/telegram"
)

const ADMIN_ID = 1152107887

func main() {
	configs.LoadConfigs()

	app := fiber.New()

	if err := database.InitConnection(); err != nil {
		panic("database connection failed.")
	}

	database.Migration()

	updates := telegram.InitBot()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	app = routes.DefineRoutes(app)

	go app.Listen(":" + configs.ServerPort)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			query := strings.Split(update.Message.Text, " ")

			if len(query) > 1 && query[1] == "otp" {
				database.DB.Where("userid = ?", update.Message.Chat.ID).Delete(&database.OTP{})
				code := uuid.NewString()
				database.DB.Create(&database.OTP{
					Code:     code,
					Userid:   update.Message.Chat.ID,
					Username: update.Message.From.UserName,
				})
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					"Click button below to confirm your account.")
				msg.ReplyToMessageID = update.Message.MessageID

				url := configs.Url + "/confirm/" + code

				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Confirm", url),
					),
				)

				telegram.Bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
				msg.ReplyToMessageID = update.Message.MessageID

				telegram.Bot.Send(msg)
			}
		case "backup":
			if update.Message.Chat.ID == ADMIN_ID {
				file := tgbotapi.FilePath("./local.db")
				telegram.Bot.Send(tgbotapi.NewDocument(update.Message.Chat.ID, file))
			}
		}

	}
}
