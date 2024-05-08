package main

import (
	"log"
	"os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitConnection() error {
	var DB_PATH = os.Getenv("DB_PATH")

	var err error
	db, err = gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})

	if err != nil {
		return err
	}

	log.Println("Database connection established successfully.")

	return nil
}

func Migration()  {
	db.AutoMigrate(&OTP{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Message{})
}