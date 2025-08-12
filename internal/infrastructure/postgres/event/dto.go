package event

import (
	"time"

	domain "billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/utils"
	"gorm.io/gorm"
)

//
// Event
//

type eventEntity struct {
	gorm.Model
	ID              int64      `gorm:"column:id"`
	Name            string     `gorm:"column:name"`
	EventType       string     `gorm:"column:event_type"`
	Status          string     `gorm:"column:status"`
	CreatedByUserID int64      `gorm:"column:created_by_user_id"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at"`
	EventDate       *time.Time `gorm:"column:event_date"`
}

func (eventEntity) TableName() string {
	return "event"
}

func eventFromDomain(d *domain.Event) *eventEntity {
	if d == nil {
		return nil
	}

	return &eventEntity{
		ID:              d.ID,
		Name:            d.Name,
		EventType:       string(d.Type),
		Status:          string(d.Status),
		CreatedByUserID: int64(utils.SafeDereference(d.CreatedBy.UserID)),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		DeletedAt:       d.DeletedAt,
		EventDate:       d.EventDate,
	}
}

func eventToDomain(e *eventEntity, createdBy domain.Member, members []domain.Member) *domain.Event {
	if e == nil {
		return nil
	}

	out := domain.Event{
		ID:        e.ID,
		Name:      e.Name,
		Members:   members,
		CreatedBy: createdBy,
		Status:    domain.Status(e.Status),
		Type:      domain.Type(e.EventType),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		EventDate: e.EventDate,
		DeletedAt: e.DeletedAt,
	}

	return &out
}

//
// Member
//

type memberEntity struct {
	gorm.Model
	ID      int64  `gorm:"column:id"`
	Name    string `gorm:"column:name"`
	EventID int64  `gorm:"column:event_id"`
	UserID  *int64 `gorm:"column:user_id"`
}

func (memberEntity) TableName() string {
	return "member"
}

func memberFromDomain(eventID int64, d *domain.Member) *memberEntity {
	if d == nil {
		return nil
	}

	out := &memberEntity{
		ID:      d.ID,
		Name:    d.Name,
		EventID: eventID,
	}

	if d.UserID != nil {
		out.UserID = utils.Ptr(int64(*d.UserID))
	}

	return out
}

func memberToDomain(e *memberEntity) *domain.Member {
	if e == nil {
		return nil
	}

	out := domain.Member{
		ID:   e.ID,
		Name: e.Name,
	}

	if e.UserID != nil {
		out.UserID = utils.Ptr(vo.UserID(*e.UserID))
	}

	return &out
}
