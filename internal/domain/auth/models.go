package auth

import (
	"time"

	"gorm.io/gorm"
)

type UserExtra struct {
	TelegramID string `json:"telegram_id"`
}

type User struct {
	gorm.Model
	ID        string
	Username  string
	FirstName string
	LastName  string
	Extra     UserExtra
}

type Session struct {
	gorm.Model
	ID       string
	UserID   string
	ExpireAt *time.Time
}
