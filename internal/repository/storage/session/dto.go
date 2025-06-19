package session

import (
	"time"

	"billsplitter-monolith/internal/domain/auth"
	"gorm.io/gorm"
)

type sessionEntity struct {
	gorm.Model
	ID       string     `gorm:"column:id"`
	UserID   string     `gorm:"column:user_id"`
	ExpireAt *time.Time `gorm:"column:expire_at"`
}

func (sessionEntity) TableName() string {
	return "sessions"
}

func fromDomain(d *auth.Session) *sessionEntity {
	if d == nil {
		return nil
	}

	return &sessionEntity{
		ID:       d.ID,
		UserID:   d.UserID,
		ExpireAt: d.ExpireAt,
	}
}

func toDomain(e *sessionEntity) *auth.Session {
	if e == nil {
		return nil
	}

	return &auth.Session{
		ID:       e.ID,
		UserID:   e.UserID,
		ExpireAt: e.ExpireAt,
	}
}
