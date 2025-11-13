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

type Participant struct {
	MemberID int64
	Amount   int64
}

type Bill struct {
	ID int64
	// Название чека
	Name string
	// CreatedBy - кто создал чек. Чеки могут создавать не только владельцы мита
	CreatedBy vo.UserID
	// Participants - юзеры в чеке
	Participants []Participant
	// PaidBy - ID участника (member), который оплатил чек
	PaidBy  int64
	EventID int64
	// TotalAmount - полная сумма чека
	TotalAmount vo.Amount
	// Currency - валюта в которой оплачивался чек
	Currency vo.CurrencyCode
	// Способ разбивки чека
	SplitTypeID vo.SplitType

	DeletedAt *time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}
