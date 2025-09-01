package meet

import (
	"time"
)

type CreateEventRq struct {
	EventName string     `json:"name"`
	Date      *time.Time `json:"date"`
	Members   []string   `json:"members"`
}

type Member struct {
	UserID   *int64 `json:"user_id"`
	Username string `json:"username"`
}

type Event struct {
	ID              int64      `json:"ID"`
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
