package telegram

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PrivacyForge/nashenas/configs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func InitBot() tgbotapi.UpdatesChannel {
	var err error
	var httpProxyString string
	var httpProxyUrl *url.URL

	flag.StringVar(&httpProxyString, "http-proxy", "", "Set custom proxy server.")
	flag.Parse()

	if len(strings.TrimSpace(httpProxyString)) > 0 {
		httpProxyUrl, err = url.Parse(httpProxyString)

		if err != nil {
			panic("Proxy Server Is Not Valid")
		}

		httpProxy := http.ProxyURL(httpProxyUrl)

		Bot, err = tgbotapi.NewBotAPIWithClient(configs.BotToken, "https://api.telegram.org/bot%s/%s", &http.Client{
			Transport: &http.Transport{
				Proxy:               httpProxy,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		})
	} else {
		Bot, err = tgbotapi.NewBotAPI(configs.BotToken)
	}

	if err != nil {
		panic("telegram bot connection failed.")
	}

	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	return updates
}
