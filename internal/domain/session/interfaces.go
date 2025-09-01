package session

import (
	"context"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type Service interface {
	Create(ctx context.Context, userID vo.UserID) (string, error)
	GetByID(ctx context.Context, sessionID string) (*Session, error)
}

type Repository interface {
	Create(ctx context.Context, session Session) (*Session, error)
	GetByID(ctx context.Context, id string) (*Session, error)
}
