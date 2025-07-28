package session

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, userID int64) (string, error)
	GetByID(ctx context.Context, sessionID string) (*Session, error)
}

type Repository interface {
	Create(ctx context.Context, session Session) (*Session, error)
	GetByID(ctx context.Context, id string) (*Session, error)
}
