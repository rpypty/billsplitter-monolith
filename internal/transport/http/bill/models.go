package bill

import "time"

type CreateBillRq struct {
	EventID      int64           `json:"event_id"`
	Name         string          `json:"name"`
	TotalAmount  int64           `json:"total_amount"`
	Currency     string          `json:"currency"`
	SplitType    string          `json:"split_type"`
	Participants []ParticipantRq `json:"participants"`
}

type ParticipantRq struct {
	MemberID int64 `json:"member_id"`
	Amount   int64 `json:"amount"`
}

type Participant struct {
	MemberID int64 `json:"member_id"`
	Amount   int64 `json:"amount"`
}

type Bill struct {
	ID              int64         `json:"id"`
	EventID         int64         `json:"event_id"`
	Name            string        `json:"name"`
	CreatedByUserID int64         `json:"created_by_user_id"`
	TotalAmount     int64         `json:"total_amount"`
	Currency        string        `json:"currency"`
	SplitType       string        `json:"split_type"`
	Participants    []Participant `json:"participants"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       *time.Time    `json:"deleted_at,omitempty"`
}
