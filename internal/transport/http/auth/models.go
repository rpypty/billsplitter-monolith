package auth

import "billsplitter-monolith/internal/domain/auth"

type LoginTelegramReq struct {
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	TelegramID int64  `json:"telegramID" validate:"required"`
}

type LoginTelegramRes struct {
	SessionID string `json:"sessionID"`
}

type MeRes struct {
	User *auth.User `json:"user"`
}
