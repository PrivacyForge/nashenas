package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

func ValidateToken(tokenString string, secret []byte) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims["id"].(float64), nil

}

const ADMIN_ID = 1152107887

func main() {
	godotenv.Load(".env")
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")

	if err := InitConnection(); err != nil {
		panic("database connection failed.")
	}

	Migration()

	updates, err := InitBot()

	if err != nil {
		panic("telegram bot connection failed.")
	}

	http.HandleFunc("/get-messages", GetMessagesHandler)
	http.HandleFunc("/send-message", SendMessageHandler)
	http.HandleFunc("/get-profile", GetProfileHandler)
	http.HandleFunc("/set-key", SetKeyHandler)
	http.HandleFunc("/set-username", SetUsernameHandler)
	http.HandleFunc("/confirm", ConfirmHandler)
	http.HandleFunc("/me", MeHandler)
	http.HandleFunc("/", HelloWorldHandler)

	go http.ListenAndServe(":"+PORT, nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			query := strings.Split(update.Message.Text, " ")

			if len(query) > 1 && query[1] == "otp" {
				db.Where("userid = ?", update.Message.Chat.ID).Delete(&OTP{})
				code := uuid.NewString()
				db.Create(&OTP{
					Code:     code,
					Userid:   update.Message.Chat.ID,
					Username: update.Message.From.UserName,
				})
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					"Click button below to confirm your account.")
				msg.ReplyToMessageID = update.Message.MessageID

				url := URL + "/confirm/" + code

				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Confirm", url),
					),
				)

				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		case "backup":
			if update.Message.Chat.ID == ADMIN_ID {
				file := tgbotapi.FilePath("./local.db")
				bot.Send(tgbotapi.NewDocument(update.Message.Chat.ID, file))
			}
		}

	}
}
