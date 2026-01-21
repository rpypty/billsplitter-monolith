package event

import (
	"time"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type UpdateEventRq struct {
	Name      *string
	Status    *Status
	EventDate *time.Time
}

type Member struct {
	ID     int64
	UserID *vo.UserID
	Name   string
}

type CreateEventRq struct {
	Name      string
	Type      Type
	EventDate *time.Time

	// CreatedBy - пользователя который создает мит - это сразу первый авторизованный пользователь
	CreatorUserID   vo.UserID
	CreatorUsername string

	// Members - Список участников
	// Работаем по системе "пикми" - при создании мита добавляем плейсхолдеры пользователей - просто имена.
	// Потом пользователь должен будет открыть мит
	// и нажать в списке мемберов на кнопку "Авторизовать себя".
	// В этот момент на бэке система подвяжет авторизованного (текущего) пользователя к этому мемберу
	Members []string
}

type Event struct {
	ID              int64
	PublicUUID      string
	Name            string
	CreatedByUserID vo.UserID
	Members         []Member // Members - участников, обогощаются на уровне UseCase
	Status          Status
	Type            Type
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EventDate       *time.Time
	DeletedAt       *time.Time
}
