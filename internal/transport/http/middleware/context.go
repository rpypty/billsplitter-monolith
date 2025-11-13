package middleware

import (
	"context"

	"billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/errors"
)

type ctxKey string

// InjectUserInSession injects session info into context .
func InjectUserInSession(ctx context.Context, sess *session.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, sess)
}

func SessionFromContext(ctx context.Context) (*session.Session, error) {
	u, ok := ctx.Value(sessionContextKey).(*session.Session)
	if !ok {
		return nil, errors.ErrFailedToGetUserFromCtx
	}
	return u, nil
}

func ExtractUserFromContext(ctx context.Context) (vo.UserID, error) {
	s, err := SessionFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return s.UserID, nil
}
