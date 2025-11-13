package testkit

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	domainsession "billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

// WithUserSession injects a session with the provided user ID into request context.
func WithUserSession(req *http.Request, userID vo.UserID) *http.Request {
	return WithSession(req, &domainsession.Session{UserID: userID})
}

// WithSession injects given session into request context (session may be nil).
func WithSession(req *http.Request, sess *domainsession.Session) *http.Request {
	return req.WithContext(middleware.InjectUserInSession(req.Context(), sess))
}

// WithRouteParams adds chi route params to request.
func WithRouteParams(req *http.Request, params map[string]string) *http.Request {
	routeCtx := chi.NewRouteContext()
	for k, v := range params {
		routeCtx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}

// NewTestLogger returns silent slog logger for tests.
func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
