package middleware

import (
	"context"
	"net/http"

	"billsplitter-monolith/internal/domain/auth"
)

type MWFunc func(next http.Handler) http.Handler

type Manager interface {
	Auth() MWFunc
}

type UserGetterSvc interface {
	GetUserBySessionID(ctx context.Context, sessionID string) (*auth.User, error)
}
