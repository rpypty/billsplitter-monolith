package payment_method

import (
	domain "billsplitter-monolith/internal/domain/payment_method"
	"gorm.io/gorm"
)

type paymentMethodEntity struct {
	gorm.Model
	ID          int64  `gorm:"column:id"`
	UserID      int64  `gorm:"column:user_id"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Recipient   string `gorm:"column:recipient"`
}

func (paymentMethodEntity) TableName() string {
	return "payment_methods"
}

func fromDomain(d *domain.PaymentMethod) *paymentMethodEntity {
	if d == nil {
		return nil
	}

	return &paymentMethodEntity{
		ID:          d.ID,
		UserID:      d.UserID,
		Name:        d.Name,
		Description: d.Description,
		Recipient:   d.Recipient,
	}
}

func toDomain(e *paymentMethodEntity) *domain.PaymentMethod {
	if e == nil {
		return nil
	}

	return &domain.PaymentMethod{
		ID:          e.ID,
		UserID:      e.UserID,
		Name:        e.Name,
		Description: e.Description,
		Recipient:   e.Recipient,
	}
}
