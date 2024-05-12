package database

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	TelegramUserid uint64 `gorm:"size: 255"`
	Username       string `gorm:"size: 255"`
	Code           string `gorm:"size: 255"`
}

type User struct {
	gorm.Model
	ID             uint64 `gorm:"primaryKey"`
	TelegramUserid uint64 `gorm:"size: 255"`
	Username       string `gorm:"size: 255"`
	PublicKey      string `gorm:"size: 255"`
	Messages       []Message
}

type Message struct {
	gorm.Model
	ID       uint64    `gorm:"primaryKey"`
	Content  string    `gorm:"not null"`
	Time     time.Time `gorm:"size: 255"`
	UserID   int64     `gorm:"not null"`
	User     User      `gorm:"foreignKey:UserID"`
	ParentID uint64    `gorm:"default:null"`
	Replies  []Message `gorm:"foreignKey:ParentID"`
}
