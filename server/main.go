package main

import (
	// "encoding/json"
	// "fmt"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt/v5"
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
	ID      int64     `gorm:"primaryKey"`
	Message string    `gorm:"size: 255"`
	UserId  int64     `gorm:"size: 255"`
	Time    time.Time `gorm:"size: 255"`
}

func ValidateToken(tokenString string, secret []byte) float64 {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims["id"].(float64)

}

func main() {
	godotenv.Load(".dev.env")
	PORT := os.Getenv("PORT")
	DB_PATH := os.Getenv("DB_PATH")
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	SECRET := os.Getenv("SECRET")

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

	http.HandleFunc("/get-messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

		if r.Method == "OPTIONS" {
			return
		}

		token := r.Header.Get("Authorization")

		id := ValidateToken(token, []byte(SECRET))

		var result []Message
		db.Where("user_id = ?", int64(id)).Find(&result)

		content, _ := json.Marshal(result)

		io.WriteString(w, string(content))
	})

	http.HandleFunc("/send-message", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

		if r.Method == "OPTIONS" {
			return
		}

		type Body struct {
			Message string `json:"message"`
		}

		id := r.URL.Query().Get("id")

		if id == "" {
			w.WriteHeader(400)
			io.WriteString(w, "The id query string is empty.")
			return
		}

		matched, _ := regexp.MatchString("^[0-9]*$", id)

		if !matched {
			w.WriteHeader(400)
			io.WriteString(w, "The id query string is invalid.")
			return
		}

		body, _ := io.ReadAll(r.Body)

		var resbody Body
		json.Unmarshal(body, &resbody)

		if resbody.Message == "" {
			w.WriteHeader(400)
			io.WriteString(w, "The message field body is empty.")
			return
		}

		userid, _ := strconv.ParseInt(id, 10, 64)

		db.Create(&Message{Message: resbody.Message, UserId: userid, Time: time.Now()})

		io.WriteString(w, "The message Sent.")
	})

	http.HandleFunc("/get-profile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

		if r.Method == "OPTIONS" {
			return
		}

		username := r.URL.Query().Get("username")

		if username == "" {
			w.WriteHeader(400)
			io.WriteString(w, "The username query string is empty.")
			return
		}

		matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", username)

		if !matched {
			w.WriteHeader(400)
			io.WriteString(w, "Username is invalid.")
			return
		}

		var result User
		db.Where("username = ?", username).Find(&result)

		if result.ID == 0 {
			w.WriteHeader(400)
			io.WriteString(w, "Not found.")
			return
		}

		content, _ := json.Marshal(result)

		io.WriteString(w, string(content))

	})

	http.HandleFunc("/set-key", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

		if r.Method == "OPTIONS" {
			return
		}

		type Body struct {
			PublicKey string `json:"public_key"`
		}
		token := r.Header.Get("Authorization")

		id := ValidateToken(token, []byte(SECRET))

		body, _ := io.ReadAll(r.Body)

		var resbody Body
		json.Unmarshal(body, &resbody)

		var result User
		db.Model(&result).Where("id = ?", id).Update("public_key", resbody.PublicKey)

		if result.PublicKey != "" {
			io.WriteString(w, "Key set successfully.")
		} else {
			w.WriteHeader(400)
			io.WriteString(w, "Error.")

		}

	})

	http.HandleFunc("/set-username", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

		if r.Method == "OPTIONS" {
			return
		}

		type Body struct {
			Username string `json:"username"`
		}
		token := r.Header.Get("Authorization")

		id := ValidateToken(token, []byte(SECRET))

		body, _ := io.ReadAll(r.Body)

		var resbody Body
		json.Unmarshal(body, &resbody)

		var result User
		db.Model(&result).Where("id = ?", id).Update("username", resbody.Username)

		io.WriteString(w, "Username set successfully.")

	})

	http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")

		if r.Method == "OPTIONS" {
			return
		}

		code := r.URL.Query().Get("code")

		var result OTP
		db.Where("code = ?", code).Find(&result)

		if result.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Code is invalid")

			return
		}

		var newUser = User{Userid: result.Userid, Username: result.Username, PublicKey: ""}
		db.Create(&newUser)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":     newUser.ID,
			"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(SECRET))

		io.WriteString(w, tokenString)

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

					url := "http://192.168.152.239:5173/confirm/" + code

					fmt.Println(url)

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

			}

		}
	}
}
