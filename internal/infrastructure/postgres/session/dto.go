package session

import (
	"time"

	domain "billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"gorm.io/gorm"
)

type sessionEntity struct {
	gorm.Model
	ID       string     `gorm:"column:id"`
	UserID   int64      `gorm:"column:user_id"`
	ExpireAt *time.Time `gorm:"column:expire_at"`
}

func (sessionEntity) TableName() string {
	return "sessions"
}

func fromDomain(d *domain.Session) *sessionEntity {
	if d == nil {
		return nil
	}

	return &sessionEntity{
		ID:       d.ID,
		UserID:   int64(d.UserID),
		ExpireAt: d.ExpireAt,
	}
}

func toDomain(e *sessionEntity) *domain.Session {
	if e == nil {
		return nil
	}

	return &domain.Session{
		ID:       e.ID,
		UserID:   vo.UserID(e.UserID),
		ExpireAt: e.ExpireAt,
	}
}
