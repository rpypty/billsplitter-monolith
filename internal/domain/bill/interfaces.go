package bill

import (
	"context"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type Service interface {
	// Create - создать новый чек в мите
	Create(ctx context.Context, bill Bill) (int64, error)

	// FetchByEventID - получить все чеки конкретного ивента
	FetchByEventID(ctx context.Context, eventID int64) ([]Bill, error)

	// Update - обновляет поля в чеке
	Update(ctx context.Context, rq UpdateBillRq) error

	// AddUsers - добавляет юзеров в чек
	AddUsers(ctx context.Context, billID int64, users []vo.UserID) error

	// RemoveUsers - удаляет юзеров из чека
	RemoveUsers(ctx context.Context, billID int64, users []vo.UserID) error

	// Delete - удаляет чек. Используется soft delete (deletedAt поле)
	Delete(ctx context.Context, billID int64) error
}

type Repository interface {
	Create(ctx context.Context, bill Bill) (int64, error)
	FetchByEventID(ctx context.Context, eventID int64) ([]Bill, error)
}
