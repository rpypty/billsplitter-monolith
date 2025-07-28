package user

import "context"

type Service interface {
	Fetch(ctx context.Context, filter FetchFilter) ([]User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByTelegramID(ctx context.Context, tgID int64) (*User, error)
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, id int64, user User) error
}

type Repository interface {
	GetByTelegramID(ctx context.Context, tgID int64) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Fetch(ctx context.Context, filter FetchFilter) ([]User, error)
	Update(ctx context.Context, id int64, user User) error
	Create(ctx context.Context, user User) (*User, error)
}
