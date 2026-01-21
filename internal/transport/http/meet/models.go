package meet

import (
	"time"
)

type CreateEventRq struct {
	EventName string     `json:"name"`
	Date      *time.Time `json:"date"`
	Members   []string   `json:"members"`
}

type SelectMemberRq struct {
	MemberID int64 `json:"member_id"`
}

type Member struct {
	ID       int64  `json:"ID"`
	UserID   *int64 `json:"user_id"`
	Username string `json:"username"`
}

type Event struct {
	ID              int64      `json:"ID"`
	PublicUUID      string     `json:"public_uuid"`
	Name            string     `json:"name"`
	CreatedByUserID int64      `json:"created_by_user_id"`
	Members         []Member   `json:"members"`
	Status          string     `json:"status"`
	Type            string     `json:"type"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	EventDate       *time.Time `json:"event_date,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type EventSummary struct {
	Balances    []Balance    `json:"balances"`
	Settlements []Settlement `json:"settlements"`
}

type Balance struct {
	MemberID   int64  `json:"member_id"`
	UserID     *int64 `json:"user_id"`
	Name       string `json:"name"`
	TotalPaid  int64  `json:"total_paid"`
	TotalShare int64  `json:"total_share"`
	Balance    int64  `json:"balance"`
}

type Settlement struct {
	FromMemberID int64 `json:"from_member_id"`
	ToMemberID   int64 `json:"to_member_id"`
	Amount       int64 `json:"amount"`
}
