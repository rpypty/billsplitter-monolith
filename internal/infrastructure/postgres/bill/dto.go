package bill

import (
	"time"

	domain "billsplitter-monolith/internal/domain/bill"
	vo "billsplitter-monolith/internal/domain/valueobject"

	"gorm.io/gorm"
)

type billEntity struct {
	gorm.Model
	ID              int64                   `gorm:"column:id"`
	EventID         int64                   `gorm:"column:event_id"`
	Name            string                  `gorm:"column:name"`
	CreatedByUserID int64                   `gorm:"column:created_by_user_id"`
	TotalAmount     int64                   `gorm:"column:total_amount"`
	Currency        string                  `gorm:"column:currency"`
	SplitType       string                  `gorm:"column:split_type"`
	CreatedAt       time.Time               `gorm:"column:created_at"`
	UpdatedAt       time.Time               `gorm:"column:updated_at"`
	DeletedAt       *time.Time              `gorm:"column:deleted_at"`
	Participants    []billParticipantEntity `gorm:"foreignKey:BillID;references:ID"`
}

func (billEntity) TableName() string {
	return "bills"
}

type billParticipantEntity struct {
	ID       int64 `gorm:"column:id"`
	BillID   int64 `gorm:"column:bill_id"`
	MemberID int64 `gorm:"column:member_id"`
	Amount   int64 `gorm:"column:amount"`
}

func (billParticipantEntity) TableName() string {
	return "bill_participants"
}

func billFromDomain(d *domain.Bill) *billEntity {
	if d == nil {
		return nil
	}

	participants := make([]billParticipantEntity, 0, len(d.Participants))

	for _, p := range d.Participants {
		participants = append(participants, billParticipantEntity{
			MemberID: p.MemberID,
			Amount:   p.Amount,
		})
	}

	return &billEntity{
		ID:              d.ID,
		EventID:         d.EventID,
		Name:            d.Name,
		CreatedByUserID: int64(d.CreatedBy),
		TotalAmount:     int64(d.TotalAmount),
		Currency:        string(d.Currency),
		SplitType:       string(d.SplitTypeID),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		DeletedAt:       d.DeletedAt,
		Participants:    participants,
	}
}

func billToDomain(e *billEntity) *domain.Bill {
	if e == nil {
		return nil
	}

	participants := make([]domain.Participant, 0, len(e.Participants))

	for _, p := range e.Participants {
		participants = append(participants, domain.Participant{
			MemberID: p.MemberID,
			Amount:   p.Amount,
		})
	}

	return &domain.Bill{
		ID:           e.ID,
		Name:         e.Name,
		CreatedBy:    vo.UserID(e.CreatedByUserID),
		Participants: participants,
		EventID:      e.EventID,
		TotalAmount:  vo.Amount(e.TotalAmount),
		Currency:     vo.CurrencyCode(e.Currency),
		SplitTypeID:  vo.SplitType(e.SplitType),
		DeletedAt:    e.DeletedAt,
		UpdatedAt:    e.UpdatedAt,
		CreatedAt:    e.CreatedAt,
	}
}
