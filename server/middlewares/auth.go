package middlewares

import (
	"strings"

	"github.com/PrivacyForge/nashenas/configs"
	"github.com/PrivacyForge/nashenas/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func OptionalBearerToken(c *fiber.Ctx) error {
	return BearerToken(c, true)
}

func RequiredBearerToken(c *fiber.Ctx) error {
	return BearerToken(c, false)
}

func BearerToken(c *fiber.Ctx, optional bool) error {
	authHeader := c.Get("Authorization")

	if optional && authHeader == "" {
		return c.Next()
	} else if authHeader == "" && !optional {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error{
			Message: "Unauthorized",
		})
	}

	// Check if the Authorization header starts with "Bearer "
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error{
			Message: "Unauthorized",
		})
	}

	token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Secret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error{
			Message: "Unauthorized",
		})
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	c.Locals("id", claims["id"])

	return c.Next()
}
