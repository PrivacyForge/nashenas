package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-jwt/jwt/v5"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world!")
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == "OPTIONS" {
		return
	}

	SECRET := os.Getenv("SECRET")

	token := r.Header.Get("Authorization")

	id, err := ValidateToken(token, []byte(SECRET))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var result User
	db.Where("id = ?", id).Find(&result)

	if result.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	content, _ := json.Marshal(result)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, string(content))

}

func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	type Response struct {
		Token     string `json:"token"`
		ID        int64  `json:"id"`
		Userid    int64  `json:"userid"`
		Username  string `json:"username"`
		PublicKey string `json:"publickey"`
	}

	if r.Method == "OPTIONS" {
		return
	}

	SECRET := os.Getenv("SECRET")

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

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if res.ID == 0 {
		var newUser = User{Userid: result.Userid, Username: result.Username, PublicKey: ""}
		db.Create(&newUser)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":     newUser.ID,
			"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(SECRET))

		var response = Response{
			Token:     tokenString,
			ID:        newUser.ID,
			Userid:    newUser.Userid,
			Username:  newUser.Username,
			PublicKey: newUser.PublicKey,
		}

		content, _ := json.Marshal(response)

		io.WriteString(w, string(content))
		return

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     res.ID,
		"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(SECRET))

	db.Where("code = ?", code).Delete(&OTP{})

	var response = Response{
		Token:     tokenString,
		ID:        res.ID,
		Userid:    res.Userid,
		Username:  res.Username,
		PublicKey: res.PublicKey,
	}

	content, _ := json.Marshal(response)

	io.WriteString(w, string(content))
}

func SetUsernameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if r.Method == "OPTIONS" {
		return
	}

	SECRET := os.Getenv("SECRET")

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

}

func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if r.Method == "OPTIONS" {
		return
	}

	SECRET := os.Getenv("SECRET")

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

}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
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
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if r.Method == "OPTIONS" {
		return
	}

	type Body struct {
		Message string `json:"message"`
	}

	URL := os.Getenv("URL")

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

	var result User
	db.Where("id = ?", id).Find(&result)

	if result.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Not found.")
		return
	}

	msg := tgbotapi.NewMessage(result.Userid, "You received a new message.")

	url := URL + "/inbox"

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Show", url),
		),
	)

	bot.Send(msg)

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
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if r.Method == "OPTIONS" {
		return
	}

	SECRET := os.Getenv("SECRET")

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
}
