package handlers

import (
	"fmt"
	"regexp"
	"strconv"
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
		ID:               result.ID,
		Username:         result.Username,
		Userid:           result.TelegramUserid,
		ReceivePublicKey: result.ReceivePublicKey,
		SendPublicKey:    result.SendPublicKey,
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
	database.DB.Where("telegram_userid = ?", result.TelegramUserid).Find(&res)

	// if user does not exist
	if res.ID == 0 {
		var newUser = database.User{
			TelegramUserid:   result.TelegramUserid,
			Username:         "",
			ReceivePublicKey: "",
			SendPublicKey:    ""}

		database.DB.Create(&newUser)

		token, _ := utils.GenerateToken(newUser.ID)
		var response = response.Confirm{
			Token:            token,
			ID:               newUser.ID,
			Userid:           newUser.TelegramUserid,
			Username:         newUser.Username,
			ReceivePublicKey: newUser.ReceivePublicKey,
			SendPublicKey:    newUser.SendPublicKey,
		}
		return c.JSON(response)
	}

	token, _ := utils.GenerateToken(res.ID)

	var response = response.Confirm{
		Token:            token,
		ID:               res.ID,
		Userid:           res.TelegramUserid,
		Username:         res.Username,
		ReceivePublicKey: res.ReceivePublicKey,
		SendPublicKey:    res.SendPublicKey,
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
		if res.ID == uint64(id) {
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

	if body.ReceivePublicKey == "" || body.SendPublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{
			Message: "Public key fields are empty",
		})
	}

	var result database.User
	database.DB.
		Model(&result).
		Where("id = ?", id).
		Update("receive_public_key", body.ReceivePublicKey).
		Update("send_public_key", body.SendPublicKey)

	return c.JSON(response.SetPublicKey{
		ReceivePublicKey: body.ReceivePublicKey,
		SendPublicKey:    body.SendPublicKey,
		Message:          "Public key set successfully",
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
		PublicKey: result.ReceivePublicKey,
	})
}

func SendMessage(c *fiber.Ctx) error {
	id := c.Locals("id")

	var isAnonymouse bool = false

	if id == nil {
		isAnonymouse = true
	}

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

	if isAnonymouse {
		database.DB.Create(&database.Message{
			Content: body.Message,
			ToID:    body.Id,
			OwnerID: body.Id,
			Time:    time.Now()})
	} else {
		database.DB.Create(&database.Message{
			Content: body.Message,
			FromID:  uint64(id.(float64)),
			ToID:    body.Id,
			OwnerID: body.Id,
			Time:    time.Now()})
	}

	// send alarm by telegram bot
	msg := tgbotapi.NewMessage(int64(result.TelegramUserid), "You received a new message.")

	url := configs.Url + "/inbox"

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Show", url),
		),
	)

	telegram.Bot.Send(msg)

	return c.JSON(response.SendMessage{Message: "The message was sent"})
}

func GetPublicKey(c *fiber.Ctx) error {
	messageId, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	var result database.Message
	database.DB.Where("id = ?", uint64(messageId)).Find(&result)

	var res database.User
	database.DB.Where("id = ?", result.FromID).Find(&res)

	return c.JSON(res.SendPublicKey)
}

func ReplayMessage(c *fiber.Ctx) error {
	id := c.Locals("id")

	var body request.ReplayMessage

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	var result database.Message
	database.DB.Where("id = ?", body.MessageId).Find(&result)

	database.DB.Create(&database.Message{
		Content:  body.Message,
		FromID:   uint64(id.(float64)),
		ToID:     result.FromID,
		OwnerID:  result.OwnerID,
		ParentID: result.ID,
		Time:     time.Now()})

	return c.SendString("Ok")
}

func GetMessages(c *fiber.Ctx) error {
	id := c.Locals("id").(float64)

	var result []database.Message

	database.DB.Where("to_id = ?", uint64(id)).Find(&result)

	var messages []response.GetMessages

	for i := 0; i < len(result); i++ {
		var owner bool = result[i].OwnerID == uint64(id)
		if result[i].ParentID != 0 {
			var res database.Message
			database.DB.Where("id = ?", result[i].ParentID).Find(&res)
			messages = append(messages, response.GetMessages{
				ID:      result[i].ID,
				Content: result[i].Content,
				Time:    result[i].Time,
				Owner:   owner,
				Quote: &response.Quote{
					ID:      res.ID,
					Content: res.Content,
				},
			})
		} else {
			messages = append(messages, response.GetMessages{
				ID:      result[i].ID,
				Content: result[i].Content,
				Time:    result[i].Time,
				Owner:   owner,
				Quote:   nil,
			})
		}

	}

	return c.JSON(messages)
}
