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
	ID               uint64    `gorm:"primaryKey"`
	TelegramUserid   uint64    `gorm:"size: 255"`
	Username         string    `gorm:"size: 255"`
	SendPublicKey    string    `gorm:"default:null"`
	ReceivePublicKey string    `gorm:"default:null"`
	SentMessages     []Message `gorm:"foreignKey:FromID; references: ID"`
	ReceivedMessages []Message `gorm:"foreignKey:ToID; references: ID"`
	OwnedMessages    []Message `gorm:"foreignKey:OwnerID; references: ID"`
}

type Message struct {
	gorm.Model
	ID       uint64    `gorm:"primaryKey"`
	Content  string    `gorm:"not null"`
	Time     time.Time `gorm:"not null"`
	FromID   uint64    `gorm:"default:null"`
	ToID     uint64    `gorm:"not null"`
	OwnerID  uint64    `gorm:"not null"`
	ParentID uint64    `gorm:"default:null"`
	Replies  []Message `gorm:"foreignKey:ParentID"`
}

// type Waitlist struct {
// 	gorm.Model
// 	ID         uint64 `gorm:"primaryKey"`
// 	TelegramId uint64 `gorm:"size: 255"`
// 	TwitterId  string `gorm:"default:null"`
// 	Accept     bool   `gorm:"size: 255"`
// }
