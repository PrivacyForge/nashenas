package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func HelloWorld(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	return c.SendString("Hello, World!" + fmt.Sprintf("%v", int(id)))
}

func GetMe(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	var result User
	db.Where("id = ?", int(id)).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Bad Request"})
	}
	return c.JSON(result)
}

func ConfirmOTP(c *fiber.Ctx) error {
	otp := c.Params("otp")

	var result OTP
	db.Where("code = ?", otp).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "OTP is invalid.",
		})
	}

	db.Where("code = ?", otp).Delete(&OTP{})

	var res User
	db.Where("userid = ?", result.Userid).Find(&res)

	// if user does not exist
	if res.ID == 0 {
		var newUser = User{Userid: result.Userid, Username: "", PublicKey: ""}
		db.Create(&newUser)

		token, _ := GenerateToken(newUser.ID)
		var response = ConfirmResponse{
			Token:     token,
			ID:        newUser.ID,
			Userid:    newUser.Userid,
			Username:  newUser.Username,
			PublicKey: newUser.PublicKey,
		}
		return c.JSON(response)
	}

	token, _ := GenerateToken(res.ID)

	var response = ConfirmResponse{
		Token:     token,
		ID:        res.ID,
		Userid:    res.Userid,
		Username:  res.Username,
		PublicKey: res.PublicKey,
	}

	return c.JSON(response)
}

func SetUsername(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	var body SetUsernameRequest

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	var result User
	db.Model(&result).Where("id = ?", int(id)).Update("username", body.Username)

	return c.JSON(SetUsernameResponse{
		Username: body.Username,
		Message:  "Username set successfully",
	})
}

// func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	SECRET := os.Getenv("SECRET")

// 	type Body struct {
// 		PublicKey string `json:"public_key"`
// 	}

// 	token := r.Header.Get("Authorization")

// 	id, err := ValidateToken(token, []byte(SECRET))

// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	body, _ := io.ReadAll(r.Body)

// 	var resbody Body
// 	json.Unmarshal(body, &resbody)

// 	var result User
// 	db.Model(&result).Where("id = ?", id).Update("public_key", resbody.PublicKey)

// 	if result.PublicKey != "" {
// 		io.WriteString(w, "Key set successfully.")
// 	} else {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "Error.")

// 	}

// }

// func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	username := r.URL.Query().Get("username")

// 	if username == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "The username query string is empty.")
// 		return
// 	}

// 	matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", username)

// 	if !matched {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "Username is invalid.")
// 		return
// 	}

// 	var result User
// 	db.Where("username = ?", username).Find(&result)

// 	if result.ID == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "Not found.")
// 		return
// 	}

// 	content, _ := json.Marshal(result)

// 	io.WriteString(w, string(content))
// }

// func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	type Body struct {
// 		Message string `json:"message"`
// 	}

// 	URL := os.Getenv("URL")

// 	id := r.URL.Query().Get("id")

// 	if id == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "The id query string is empty.")
// 		return
// 	}

// 	matched, _ := regexp.MatchString("^[0-9]*$", id)

// 	if !matched {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "The id query string is invalid.")
// 		return
// 	}

// 	var result User
// 	db.Where("id = ?", id).Find(&result)

// 	if result.ID == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "Not found.")
// 		return
// 	}

// 	msg := tgbotapi.NewMessage(result.Userid, "You received a new message.")

// 	url := URL + "/inbox"

// 	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonURL("Show", url),
// 		),
// 	)

// 	bot.Send(msg)

// 	body, _ := io.ReadAll(r.Body)

// 	var resbody Body
// 	json.Unmarshal(body, &resbody)

// 	if resbody.Message == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		io.WriteString(w, "The message field body is empty.")
// 		return
// 	}

// 	userid, _ := strconv.ParseInt(id, 10, 64)

// 	db.Create(&Message{Message: resbody.Message, UserId: userid, Time: time.Now()})

// 	io.WriteString(w, "The message Sent.")
// }

// func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	SECRET := os.Getenv("SECRET")

// 	token := r.Header.Get("Authorization")

// 	id, err := ValidateToken(token, []byte(SECRET))

// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	var result []Message
// 	db.Where("user_id = ?", int64(id)).Find(&result)

// 	content, _ := json.Marshal(result)

// 	io.WriteString(w, string(content))
// }
