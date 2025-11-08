package middleware

import (
	"context"

	"billsplitter-monolith/internal/domain/session"
)

// ContextWithSessionForTest injects session info into context for handler tests.
func ContextWithSessionForTest(ctx context.Context, sess *session.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, sess)
}
