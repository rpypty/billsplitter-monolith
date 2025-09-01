package event

import (
	"context"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type Service interface {
	// Create - создает ивент, возвращает айди созданного ивента
	Create(ctx context.Context, event CreateEventRq) (int64, error)

	// Update - обновляет некоторые поля в ивенте
	Update(ctx context.Context, rq UpdateEventRq) error

	// AddUsers - добавляет пачку юзеров к ивенту
	AddUsers(ctx context.Context, eventID int64, userID []vo.UserID) error

	// RemoveUsers - удаляет несколько юзеров из ивента.
	// TODO: Надо предусмотреть валидацию перед удалением - пользователя не должно быть в чеках
	RemoveUsers(ctx context.Context, eventID int64, userID []vo.UserID) error

	// Delete - удаляет ивент. Используется soft delete (deletedAt поле)
	Delete(ctx context.Context, billID int64) error

	FetchByUserID(ctx context.Context, userID vo.UserID) ([]Event, error)

	GetByID(ctx context.Context, eventID int64) (*Event, error)
}

type Repository interface {
	Create(ctx context.Context, event Event) (int64, error)

	AddMembers(ctx context.Context, eventID int64, members []Member) error

	FetchByUserID(ctx context.Context, userID vo.UserID) ([]Event, error)

	GetByID(ctx context.Context, eventID int64) (*Event, error)
}
