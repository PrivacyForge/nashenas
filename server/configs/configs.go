package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var DatabasePath string
var Url string
var ServerPort string
var Secret string
var BotToken string
var AdminId int64

func LoadConfigs() {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env variables file!")
	}

	ServerPort = os.Getenv("SERVER_PORT")
	DatabasePath = os.Getenv("DATABASE_PATH")
	BotToken = os.Getenv("BOT_TOKEN")
	Secret = os.Getenv("SECRET")
	Url = os.Getenv("URL")
	AdminId, _ = strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
}
