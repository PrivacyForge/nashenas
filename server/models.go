package main

import (
	"time"

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
