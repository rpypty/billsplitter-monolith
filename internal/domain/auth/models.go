package auth

import (
	"time"
)

type UserExtra struct {
	TelegramID int64 `json:"telegram_id,omitempty"`
}

type User struct {
	ID        string
	Username  string
	FirstName string
	LastName  string
	Extra     UserExtra
}

type Session struct {
	ID       string
	UserID   string
	ExpireAt *time.Time
}
