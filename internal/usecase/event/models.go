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
