package session

import "time"

type Session struct {
	ID       string
	UserID   int64
	ExpireAt *time.Time
}
