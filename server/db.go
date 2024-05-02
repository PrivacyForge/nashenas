package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	Userid   int64  `gorm:"size: 255"`
	Username string `gorm:"size: 255"`
	Code     string `gorm:"size: 255"`
}

type User struct {
	gorm.Model
	ID        int64  `gorm:"primaryKey"`
	Userid    int64  `gorm:"size: 255"`
	Username  string `gorm:"size: 255"`
	PublicKey string `gorm:"size: 255"`
}

type Message struct {
	gorm.Model
	ID      int64     `gorm:"primaryKey"`
	Message string    `gorm:"size: 255"`
	UserId  int64     `gorm:"size: 255"`
	Time    time.Time `gorm:"size: 255"`
}

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