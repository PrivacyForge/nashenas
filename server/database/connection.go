package database

import (
	"log"

	"github.com/PrivacyForge/nashenas/configs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitConnection() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(configs.DatabasePath), &gorm.Config{})

	if err != nil {
		return err
	}

	log.Println("Database connection established successfully.")

	return nil
}

func Migration() {
	DB.AutoMigrate(&OTP{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Message{})
}
