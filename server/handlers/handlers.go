package handlers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PrivacyForge/nashenas/database"
	"github.com/PrivacyForge/nashenas/redis"
	"github.com/PrivacyForge/nashenas/request"
	"github.com/PrivacyForge/nashenas/response"
	"github.com/PrivacyForge/nashenas/utils"
	"github.com/gofiber/fiber/v2"
)

func HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func GetMe(c *fiber.Ctx) error {
	userid := c.Locals("userid").(int64)

	var result database.User
	database.DB.Where("userid = ?", userid).Find(&result)

	if result.Userid == 0 {
		var newUser = database.User{
			Userid:        uint64(userid),
			Username:      "",
			PublicKey:     "",
			PublicKeyHash: "",
		}

		database.DB.Create(&newUser)
		return c.JSON(response.GetMe{
			Username:      newUser.Username,
			Userid:        uint64(userid),
			PublicKey:     newUser.PublicKey,
			PublicKeyHash: newUser.PublicKeyHash,
		})
	}

	return c.JSON(response.GetMe{
		Username:      result.Username,
		Userid:        result.Userid,
		PublicKey:     result.PublicKey,
		PublicKeyHash: result.PublicKeyHash,
	})
}

func SetUsername(c *fiber.Ctx) error {
	userid := c.Locals("userid").(int64)

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

	database.DB.Model(&database.User{}).
		Where("username = ?", strings.ToLower(body.Username)).
		Find(&res)

	if res.Userid != 0 {
		if res.Userid == uint64(userid) {
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
	database.DB.Model(&result).Where("userid = ?", userid).Update("username", strings.ToLower(body.Username))

	return c.JSON(response.SetUsername{
		Username: strings.ToLower(body.Username),
		Message:  "Username set successfully",
	})
}

func SetPublicKey(c *fiber.Ctx) error {
	userid := c.Locals("userid").(int64)

	var body request.SetPublicKey
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if body.PublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{
			Message: "Public key fields are empty",
		})
	}

	var result database.User
	database.DB.
		Model(&result).
		Where("userid = ?", userid).
		Update("public_key", body.PublicKey).
		Update("public_key_hash", utils.GenerateSHA1(body.PublicKey))

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

	if result.Userid == 0 {
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
	userid := c.Locals("userid")

	var body request.SendMessage

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if len(body.Message) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error{Message: "Bad Request"})
	}

	var targetUser database.User
	database.DB.Where("id = ?", body.Id).Find(&targetUser)

	if targetUser.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.Error{Message: "User not found"})
	}

	var sourceUser database.User
	database.DB.Where("userid = ?", userid).Find(&sourceUser)

	var session = database.Session{
		Key:  body.SessionKey,
		Time: time.Now()}

	database.DB.Create(&session)

	database.DB.Create(&database.Message{
		Content:   body.Message,
		FromID:    sourceUser.ID,
		SessionID: session.ID,
		ToID:      targetUser.ID,
		OwnerID:   targetUser.ID,
		Time:      time.Now()})

	// push notification from telegram bot to user account
	err := redis.Client.Publish("message", fmt.Sprint(targetUser.Userid))
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	return c.JSON(response.SendMessage{Message: "The message was sent", SessionID: session.ID})
}

func GetPublicKey(c *fiber.Ctx) error {
	messageId, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	var result database.Message
	database.DB.Where("id = ?", messageId).Find(&result)

	var res database.User
	database.DB.Where("id = ?", result.FromID).Find(&res)

	if result.OwnerID == result.FromID {
		return c.JSON(res.PublicKey)
	}

	return c.JSON(res.PublicKey)
}

func ReplayMessage(c *fiber.Ctx) error {
	userid := c.Locals("userid")

	var body request.ReplayMessage

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	var result database.Message
	database.DB.Where("id = ?", body.MessageId).Find(&result)

	var user database.User
	database.DB.Where("userid = ?", userid).Find(&user)

	database.DB.Create(&database.Message{
		Content:  body.Message,
		FromID:   user.ID,
		ToID:     result.FromID,
		OwnerID:  result.OwnerID,
		ParentID: result.ID,
		Time:     time.Now()})

	var targetUser database.User
	database.DB.Where("id = ?", result.FromID).Find(&targetUser)

	err := redis.Client.Publish("message", fmt.Sprint(targetUser.Userid))
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	return c.SendString("Ok")
}

func GetMessages(c *fiber.Ctx) error {
	userid := c.Locals("userid").(int64)

	var user database.User
	database.DB.Where("userid = ?", userid).Find(&user)

	var result []database.Message
	database.DB.Where("to_id = ?", user.ID).Find(&result)

	messages := []response.GetMessages{}

	for i := 0; i < len(result); i++ {
		var owner bool = result[i].OwnerID == uint64(user.ID)

		var sourceUser database.User
		database.DB.Where("id = ?", result[i].FromID).Find(&sourceUser)
		if result[i].ParentID != 0 {
			var res database.Message
			database.DB.Where("id = ?", result[i].ParentID).Find(&res)
			var quotedUser database.User
			database.DB.Where("id = ?", res.ToID).Find(&quotedUser)
			if owner {
				messages = append(messages, response.GetMessages{
					ID:        result[i].ID,
					Content:   result[i].Content,
					Time:      result[i].Time,
					Owner:     owner, // true
					CanReplay: true,
					Quote: &response.Quote{
						ID:      res.ID,
						Content: res.Content,
					},
				})
			} else {
				messages = append(messages, response.GetMessages{
					ID:        result[i].ID,
					Content:   result[i].Content,
					Time:      result[i].Time,
					SessionID: result[i].SessionID,
					Owner:     owner, // false
					CanReplay: true,
					Quote: &response.Quote{
						ID:      res.ID,
						Content: res.Content,
					},
				})
			}
		} else {
			var session database.Session
			database.DB.Where("id = ?", result[i].SessionID).Find(&session)

			messages = append(messages, response.GetMessages{
				ID:        result[i].ID,
				SessionID: result[i].SessionID,
				Time:      result[i].Time,
				Owner:     owner,
				Quote:     nil,
				Content:   result[i].Content,
				CanReplay: true,
			})

			// database.DB.Delete(&database.Session{}, result[i].SessionID)
		}

	}

	return c.JSON(messages)
}
