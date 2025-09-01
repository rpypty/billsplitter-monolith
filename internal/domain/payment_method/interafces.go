package payment_method

import (
	"context"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type Service interface {
	// FetchByUserID - возвращает список платежных методов для юзера
	FetchByUserID(ctx context.Context, userID vo.UserID) ([]PaymentMethod, error)
	Create(ctx context.Context, paymentMethod PaymentMethod) (PaymentMethod, error)
	Update(ctx context.Context, id int64, paymentMethod PaymentMethod) error
	Delete(ctx context.Context, id int64) error
}

type Repository interface {
	GetByID(ctx context.Context, id int64) (PaymentMethod, error)
	FetchByUserID(ctx context.Context, userID vo.UserID) ([]PaymentMethod, error)
	Create(ctx context.Context, method PaymentMethod) (PaymentMethod, error)
	Update(ctx context.Context, id int64, method PaymentMethod) error
	Delete(ctx context.Context, id int64) error
}
