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

type Event struct {
	ID           int
	Name         string
	Participants []vo.UserID // Participants - айди участников, обогощаются на уровне UseCase
	CreatedBy    vo.UserID
	Status       Status
	Type         Type
	CreatedAt    time.Time
	UpdatedAt    time.Time
	EventDate    *time.Time
	DeletedAt    *time.Time
}
