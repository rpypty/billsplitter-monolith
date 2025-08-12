package meet

import (
	"time"
)

type CreateEventRq struct {
	EventName string     `json:"name"`
	Date      *time.Time `json:"date"`
	Members   []string   `json:"members"`
}
