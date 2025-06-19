package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/cfg"
	hu "billsplitter-monolith/internal/utils/http"
)

type middlewareManagerImpl struct {
	userGetter UserGetterSvc

	logger *slog.Logger
}

func NewMiddlewareManager(getter UserGetterSvc, logger *slog.Logger) Manager {
	return &middlewareManagerImpl{
		userGetter: getter,
		logger:     logger,
	}
}

func (mw *middlewareManagerImpl) Auth() MWFunc {
	l := mw.l().With("method", "Auth")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get(xSessionID)
			if sessionID == "" {
				hu.RespondErrWithStatus(w, http.StatusUnauthorized, Unauthorized)
				return
			}

			user, err := mw.userGetter.GetUserBySessionID(r.Context(), sessionID)
			if err != nil {
				l.Error("AuthorizeMiddleware: failed to get user from session", "error", err)

				if cfg.IsDebug() {
					// показываем подробную ошибку в дебаг моде
					hu.RespondErrWithStatusf(
						w,
						http.StatusUnauthorized,
						"AuthorizeMiddleware: failed to get user from session: %s",
						err,
					)
					return
				}

				hu.RespondErrWithStatus(w, http.StatusUnauthorized, Unauthorized)
				return
			}

			if user == nil {
				hu.RespondErrWithStatus(w, http.StatusUnauthorized, Unauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *middlewareManagerImpl) l() *slog.Logger {
	return mw.logger.WithGroup("middleware-manager")
}
