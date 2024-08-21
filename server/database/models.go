package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID               uint64    `gorm:"primaryKey"`
	Userid           uint64    `gorm:"size: 255"`
	Username         string    `gorm:"size: 255"`
	PublicKey        string    `gorm:"default:null"`
	PublicKeyHash    string    `gorm:"default:null"`
	SentMessages     []Message `gorm:"foreignKey:FromID; references: ID"`
	ReceivedMessages []Message `gorm:"foreignKey:ToID; references: ID"`
	OwnedMessages    []Message `gorm:"foreignKey:OwnerID; references: ID"`
}

type Message struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey"`
	Content   string    `gorm:"not null"`
	SessionID uint64    `gorm:"not null"`
	Time      time.Time `gorm:"not null"`
	FromID    uint64    `gorm:"default:null"`
	ToID      uint64    `gorm:"not null"`
	OwnerID   uint64    `gorm:"not null"`
	ParentID  uint64    `gorm:"default:null"`
	Replies   []Message `gorm:"foreignKey:ParentID"`
}

type Session struct {
	gorm.Model
	ID       uint64    `gorm:"primaryKey"`
	Key      string    `gorm:"not null"`
	Time     time.Time `gorm:"not null"`
	Sessions []Message `gorm:"foreignKey:SessionID; references:ID"`
}
