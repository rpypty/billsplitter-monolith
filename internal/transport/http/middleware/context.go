package middleware

import (
	"context"

	"billsplitter-monolith/internal/domain/session"
	"billsplitter-monolith/internal/errors"
)

type ctxKey string

func SessionFromContext(ctx context.Context) (*session.Session, error) {
	u, ok := ctx.Value(sessionContextKey).(*session.Session)
	if !ok {
		return nil, errors.ErrFailedToGetUserFromCtx
	}
	return u, nil
}
