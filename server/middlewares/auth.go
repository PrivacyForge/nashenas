package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PrivacyForge/nashenas/configs"
	"github.com/PrivacyForge/nashenas/response"
	"github.com/PrivacyForge/nashenas/utils"
	"github.com/gofiber/fiber/v2"
)

func sign(payload, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}

func ValidateInitData(initData, token string, expIn time.Duration) error {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return err
	}

	var (
		authDate time.Time
		hash     string
		pairs    = make([]string, 0, len(q))
	)

	for k, v := range q {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err == nil {
				authDate = time.Unix(int64(i), 0)
			}
		}
		pairs = append(pairs, k+"="+v[0])
	}

	if hash == "" {
		return errors.New("hash is empty")
	}

	if expIn > 0 {
		if authDate.IsZero() {
			return errors.New("expired")
		}

		if authDate.Add(expIn).Before(time.Now()) {
			return errors.New("expired")
		}
	}

	sort.Strings(pairs)

	if sign(strings.Join(pairs, "\n"), token) != hash {
		return errors.New("sign error")
	}
	return nil
}

func AuthMiddleware(c *fiber.Ctx) error {
	token := configs.BotToken
	initData := fmt.Sprint(c.Get("Authorization"))
	expIn := 168 * time.Hour

	err := ValidateInitData(initData, token, expIn)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error{Message: "invalid telegram initdata"})
	}

	var init, _ = utils.Parse(initData)
	c.Locals("userid", init.User.ID)

	return c.Next()
}
