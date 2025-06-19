package auth

import "billsplitter-monolith/internal/domain/auth"

// LoginTelegramReq содержит данные пользователя, полученные из Telegram
type LoginTelegramReq struct {
	Username   string `json:"username" example:"johndoe"`
	FirstName  string `json:"firstName" example:"John"`
	LastName   string `json:"lastName" example:"Doe"`
	TelegramID int64  `json:"telegramID" validate:"required" example:"123456789"`
}

// LoginTelegramRes содержит sessionID, выданный после успешного входа
type LoginTelegramRes struct {
	SessionID string `json:"sessionID" example:"b42b0a8e-0d1f-4c3d-939f-85fbbdc9be62"`
}

// MeRes содержит данные авторизованного пользователя
type MeRes struct {
	User *auth.User `json:"user"`
}
