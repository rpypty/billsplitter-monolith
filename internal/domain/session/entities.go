package session

import (
	"time"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type Session struct {
	ID       string
	UserID   vo.UserID
	ExpireAt *time.Time
}
