package handlers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PrivacyForge/nashenas/configs"
	"github.com/PrivacyForge/nashenas/database"
	"github.com/PrivacyForge/nashenas/request"
	"github.com/PrivacyForge/nashenas/response"
	"github.com/PrivacyForge/nashenas/telegram"
	"github.com/PrivacyForge/nashenas/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

func HelloWorld(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	return c.SendString("Hello, World!" + fmt.Sprintf("%v", int(id)))
}

func GetMe(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	var result database.User
	database.DB.Where("id = ?", int(id)).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{Message: "Bad Request"})
	}
	return c.JSON(response.GetMe{
		ID:        result.ID,
		Username:  result.Username,
		Userid:    result.Userid,
		PublicKey: result.PublicKey,
	})
}

func ConfirmOTP(c *fiber.Ctx) error {
	otp := c.Params("otp")

	var result database.OTP
	database.DB.Where("code = ?", otp).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{
			Message: "OTP is invalid.",
		})
	}

	database.DB.Where("code = ?", otp).Delete(&database.OTP{})

	var res database.User
	database.DB.Where("userid = ?", result.Userid).Find(&res)

	// if user does not exist
	if res.ID == 0 {
		var newUser = database.User{Userid: result.Userid, Username: "", PublicKey: ""}
		database.DB.Create(&newUser)

		token, _ := utils.GenerateToken(newUser.ID)
		var response = response.Confirm{
			Token:     token,
			ID:        newUser.ID,
			Userid:    newUser.Userid,
			Username:  newUser.Username,
			PublicKey: newUser.PublicKey,
		}
		return c.JSON(response)
	}

	token, _ := utils.GenerateToken(res.ID)

	var response = response.Confirm{
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

	var body request.SetUsername

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", body.Username)

	if !matched {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{
			Message: "Username is invalid",
		})
	}

	var res database.User
	database.DB.Model(&database.User{}).Where("username = ?", body.Username).Find(&res)

	if res.ID != 0 {
		if res.ID == int64(id) {
			return c.Status(fiber.StatusOK).JSON(response.Error{
				Message: "This username was already set for you",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error{
				Message: "Username is duplicate",
			})
		}
	}

	var result database.User
	database.DB.Model(&result).Where("id = ?", int(id)).Update("username", body.Username)

	return c.JSON(response.SetUsername{
		Username: body.Username,
		Message:  "Username set successfully",
	})
}

func SetPublicKey(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)
	var body request.SetPublicKey

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if body.PublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{
			Message: "Public key is empty",
		})
	}

	var result database.User
	database.DB.Model(&result).Where("id = ?", id).Update("public_key", body.PublicKey)

	return c.JSON(response.SetPublicKey{
		PublicKey: body.PublicKey,
		Message:   "Public key set successfully",
	})
}

func GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	matched, _ := regexp.MatchString("^[a-zA-Z]{1}[a-zA-Z0-9]{4,}$", username)

	if !matched {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{
			Message: "Username is invalid",
		})
	}

	var result database.User
	database.DB.Where("username = ?", username).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.Error{
			Message: "User Not found",
		})
	}

	return c.JSON(response.GetProfile{
		ID:        result.ID,
		Username:  result.Username,
		PublicKey: result.PublicKey,
	})
}

func SendMessage(c *fiber.Ctx) error {
	var body request.SendMessage

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	matched, _ := regexp.MatchString("^[0-9]*$", fmt.Sprintf("%v", body.Id))

	if !matched || len(body.Message) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{Message: "Bad Request"})
	}

	var result database.User
	database.DB.Where("id = ?", body.Id).Find(&result)

	if result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.Error{Message: "User not found"})
	}

	database.DB.Create(&database.Message{Message: body.Message, UserId: body.Id, Time: time.Now()})

	// send alarm by telegram bot
	msg := tgbotapi.NewMessage(result.Userid, "You received a new message.")

	url := configs.Url + "/inbox"

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Show", url),
		),
	)

	telegram.Bot.Send(msg)

	return c.JSON(response.SendMessage{Message: "The message was sent"})
}

func GetMessages(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	var result []database.Message
	database.DB.Where("user_id = ?", int64(id)).Find(&result)

	return c.JSON(result)
}
