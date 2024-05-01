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
	DB_PATH := os.Getenv("DB_PATH")
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	SECRET := os.Getenv("SECRET")
	URL := os.Getenv("URL")

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

		id, err := ValidateToken(token, []byte(SECRET))

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "The id query string is empty.")
			return
		}

		matched, _ := regexp.MatchString("^[0-9]*$", id)

		if !matched {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "The id query string is invalid.")
			return
		}

		body, _ := io.ReadAll(r.Body)

		var resbody Body
		json.Unmarshal(body, &resbody)

		if resbody.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "The username query string is empty.")
			return
		}

		matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", username)

		if !matched {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Username is invalid.")
			return
		}

		var result User
		db.Where("username = ?", username).Find(&result)

		if result.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
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

		id, err := ValidateToken(token, []byte(SECRET))

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, _ := io.ReadAll(r.Body)

		var resbody Body
		json.Unmarshal(body, &resbody)

		var result User
		db.Model(&result).Where("id = ?", id).Update("public_key", resbody.PublicKey)

		if result.PublicKey != "" {
			io.WriteString(w, "Key set successfully.")
		} else {
			w.WriteHeader(http.StatusBadRequest)
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

		id, err := ValidateToken(token, []byte(SECRET))

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		code := r.URL.Query().Get("code")

		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "The code query string is empty.")
			return
		}

		var result OTP
		db.Where("code = ?", code).Find(&result)

		if result.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Code is invalid")

			return
		}

		var res User
		db.Where("userid = ?", result.Userid).Find(&res)

		fmt.Println(result.Userid, res.ID)

		if res.ID == 0 {
			var newUser = User{Userid: result.Userid, Username: result.Username, PublicKey: ""}
			db.Create(&newUser)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":     newUser.ID,
				"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(SECRET))

			io.WriteString(w, tokenString)
			return

		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":     res.ID,
			"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(SECRET))

		db.Where("code = ?", code).Delete(&OTP{})

		io.WriteString(w, tokenString)

	})

	http.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		token := r.Header.Get("Authorization")

		id, err := ValidateToken(token, []byte(SECRET))

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var result User
		db.Where("id = ?", id).Find(&result)

		content, _ := json.Marshal(result)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, string(content))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world!")
	})

	go http.ListenAndServe(":"+PORT, nil)

	for update := range updates {
		if update.Message != nil { // If we got a message

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
}
