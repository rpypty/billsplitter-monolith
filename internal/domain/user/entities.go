package user

import vo "billsplitter-monolith/internal/domain/valueobject"

type User struct {
	ID        vo.UserID
	Username  string
	FirstName string
	LastName  string
	Extra     ExtraInfo
}

type FetchFilter struct {
	BillID    *int64
	MeetingID *int64
	UserIDs   []int64
}

type ExtraInfo struct {
	TelegramID int64 `json:"telegram_id,omitempty"`
}
