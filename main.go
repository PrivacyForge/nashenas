package main

import (
	// "encoding/json"
	// "fmt"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	Userid   int64  `gorm:"size: 255"`
	Username string `gorm:"size: 255"`
	Code     string `gorm:"size: 255"`
}

type User struct {
	gorm.Model
	ID        int64  `gorm:"primaryKey"`
	Userid    int64  `gorm:"size: 255"`
	Username  string `gorm:"size: 255"`
	PublicKey string `gorm:"size: 255"`
}

type Message struct {
	gorm.Model
	ID      int64  `gorm:"primaryKey"`
	Message string `gorm:"size: 255"`
	UserId  int64  `gorm:"size: 255"`
}

type SetKeyBody struct {
	PublicKey string `json:"public_key"`
}

type SendMessageBody struct {
	Message string `json:"message"`
}

func main() {
	godotenv.Load(".dev.env")
	PORT := os.Getenv("PORT")
	DB_PATH := os.Getenv("DB_PATH")
	BOT_TOKEN := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&OTP{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Message{})

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	http.HandleFunc("GET /get-messages", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		var result []Message
		db.Where("user_id = ?", id).Find(&result)

		fmt.Println(result)

	})

	http.HandleFunc("POST /send-message", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		body, _ := io.ReadAll(r.Body)
		// message

		var resbody SendMessageBody
		json.Unmarshal(body, &resbody)

		userid, _ := strconv.ParseInt(id, 10, 64)
		db.Create(&Message{Message: resbody.Message, UserId: userid})

		io.WriteString(w, "The message Sent.")
	})

	http.HandleFunc("GET /get-profile", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		var result User
		db.Where("id = ?", id).Find(&result)

		if result.ID == 0 {
			w.WriteHeader(400)
			io.WriteString(w, "Not found.")
			return
		}

		content, _ := json.Marshal(result)

		io.WriteString(w, string(content))

	})

	http.HandleFunc("POST /set-key", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		body, _ := io.ReadAll(r.Body)

		var resbody SetKeyBody
		json.Unmarshal(body, &resbody)

		fmt.Println(resbody.PublicKey)

		var result User
		db.Model(&result).Where("id = ?", id).Update("public_key", resbody.PublicKey)

		if result.PublicKey != "" {
			io.WriteString(w, "Key set successfully.")
		} else {
			w.WriteHeader(400)
			io.WriteString(w, "Error.")

		}

	})

	http.HandleFunc("GET /sign-up", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		var result OTP
		db.Where("code = ?", code).Find(&result)

		if result.ID != 0 {
			db.Create(&User{Userid: result.Userid, Username: result.Username, PublicKey: ""})
			io.WriteString(w, "The user has been generated.")
		} else {
			w.WriteHeader(400)
			io.WriteString(w, "Code is invalid")
		}

	})

	go http.ListenAndServe(":"+PORT, nil)

	for update := range updates {
		if update.Message != nil { // If we got a message

			switch update.Message.Command() {
			case "start":
				query := strings.Split(update.Message.Text, " ")

				fmt.Println(query)
				if len(query) > 1 && query[1] == "otp" {
					code := uuid.NewString()
					db.Create(&OTP{
						Code:     code,
						Userid:   update.Message.Chat.ID,
						Username: update.Message.From.UserName,
					})
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, code)
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}

			}

		}
	}
}
