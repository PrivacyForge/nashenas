package utils

import (
	"time"

	"github.com/PrivacyForge/nashenas/configs"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     id,
		"expire": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(configs.Secret))

	return tokenString, err
}
