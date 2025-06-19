package auth

import (
	"context"
	"time"
)

type Service interface {
	GetUserBySessionID(ctx context.Context, sessionID string) (*User, error)
	CreateSession(ctx context.Context, s *Session) (string, error)
	CreateOrGetUserByTgID(ctx context.Context, tgID int64, data *User) (*User, error)
}

type UserStorage interface {
	GetByTelegramID(ctx context.Context, tgID int64) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type SessionCache interface {
	Set(ctx context.Context, sessionID string, value *Session, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (*Session, error)
}

type SessionStorage interface {
	Create(ctx context.Context, session *Session) error
	Get(ctx context.Context, id string) (*Session, error)
}
