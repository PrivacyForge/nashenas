package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", body.Username)

	if !matched {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Username is invalid",
		})
	}

	var res User
	db.Model(&User{}).Where("username = ?", body.Username).Find(&res)

	if res.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Username is duplicate",
		})
	}

	var result User
	db.Model(&result).Where("id = ?", int(id)).Update("username", body.Username)

	return c.JSON(SetUsernameResponse{
		Username: body.Username,
		Message:  "Username set successfully",
	})
}

func SetPublicKey(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)
	var body SetPublicKeyRequest

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if body.PublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Public key is empty",
		})
	}

	var result User
	db.Model(&result).Where("id = ?", id).Update("public_key", body.PublicKey)

	return c.JSON(SetPublicKeyResponse{
		PublicKey: body.PublicKey,
		Message:   "Public key set successfully",
	})
}

func GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", username)

	if !matched {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: "Username is invalid",
		})
	}

	var result User
	db.Where("username = ?", username).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Message: "User Not found",
		})
	}

	return c.JSON(GetProfileResponse{
		ID:        result.ID,
		Username:  result.Username,
		PublicKey: result.PublicKey,
	})
}

func SendMessage(c *fiber.Ctx) error {
	URL := os.Getenv("URL")

	var body SendMessageRequest

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	matched, _ := regexp.MatchString("^[0-9]*$", fmt.Sprintf("%v", body.Id))

	if !matched || len(body.Message) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Bad Request"})
	}

	var result User
	db.Where("id = ?", body.Id).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Message: "User not found"})
	}

	db.Create(&Message{Message: body.Message, UserId: body.Id, Time: time.Now()})

	// send alarm by telegram bot
	msg := tgbotapi.NewMessage(result.Userid, "You received a new message.")

	url := URL + "/inbox"

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Show", url),
		),
	)

	bot.Send(msg)

	return c.JSON(SendMessageResponse{Message: "The message was sent"})
}

func GetMessages(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	var result []Message
	db.Where("user_id = ?", int64(id)).Find(&result)

	return c.JSON(result)
}
