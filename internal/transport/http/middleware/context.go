package middleware

import (
	"context"

	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/errors"
)

type ctxKey string

func UserFromContext(ctx context.Context) (*auth.User, error) {
	u, ok := ctx.Value(userContextKey).(*auth.User)
	if !ok {
		return nil, errors.ErrFailedToGetUserFromCtx
	}
	return u, nil
}
