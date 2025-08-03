package bill

import (
	"time"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type UpdateBillRq struct {
	Name        *string
	TotalAmount *vo.Amount
	Currency    *vo.CurrencyCode
}

type Bill struct {
	ID   int64
	Name string
	// CreatedBy - кто создал чек. Чеки могут создавать не только владельцы мита
	CreatedBy vo.UserID
	// Participants - юзеры в чеке
	Participants []vo.UserID
	EventID      int64
	// TotalAmount - полная сумма чека
	TotalAmount vo.Amount
	// Currency - валюта в которой оплачивался чек
	Currency vo.CurrencyCode

	DeletedAt *time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}
