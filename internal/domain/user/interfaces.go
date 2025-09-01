package user

import (
	"context"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type Service interface {
	Fetch(ctx context.Context, filter FetchFilter) ([]User, error)
	GetByID(ctx context.Context, id vo.UserID) (*User, error)
	GetByTelegramID(ctx context.Context, tgID int64) (*User, error)
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, id vo.UserID, user User) error
}

type Repository interface {
	GetByTelegramID(ctx context.Context, tgID int64) (*User, error)
	GetByID(ctx context.Context, id vo.UserID) (*User, error)
	Fetch(ctx context.Context, filter FetchFilter) ([]User, error)
	Update(ctx context.Context, id vo.UserID, user User) error
	Create(ctx context.Context, user User) (*User, error)
}
