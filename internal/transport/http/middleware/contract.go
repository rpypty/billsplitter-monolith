package middleware

import (
	"context"
	"net/http"

	"billsplitter-monolith/internal/domain/session"
)

type MWFunc func(next http.Handler) http.Handler

type Manager interface {
	Auth() MWFunc
}

type SessionGetterSvc interface {
	GetByID(ctx context.Context, sessionID string) (*session.Session, error)
}
