package event

import (
	"time"

	vo "billsplitter-monolith/internal/domain/valueobject"
)

type CreateMeetRq struct {
	EventName       string
	Date            *time.Time
	CreatedByUserID vo.UserID
	Members         []string
}

type Event struct {
	EventName       string
	Date            *time.Time
	CreatedByUserID vo.UserID
	Members         []string
}

type EventSummary struct {
	Balances    []Balance
	Settlements []Settlement
}

type Balance struct {
	MemberID   int64
	UserID     *vo.UserID
	Name       string
	TotalPaid  int64
	TotalShare int64
	Balance    int64
}

type Settlement struct {
	FromMemberID int64
	ToMemberID   int64
	Amount       int64
}
