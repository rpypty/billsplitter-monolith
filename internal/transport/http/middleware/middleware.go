package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"billsplitter-monolith/internal/cfg"
	hu "billsplitter-monolith/internal/utils/http"
)

type middlewareManagerImpl struct {
	userGetter SessionGetterSvc

	logger *slog.Logger
}

func NewMiddlewareManager(getter SessionGetterSvc, logger *slog.Logger) Manager {
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

			sessionInfo, err := mw.userGetter.GetByID(r.Context(), sessionID)
			if err != nil {
				l.Error("AuthorizeMiddleware: failed to get session", "error", err)

				if cfg.IsDebug() {
					// показываем подробную ошибку в дебаг моде
					hu.RespondErrWithStatusf(
						w,
						http.StatusUnauthorized,
						"AuthorizeMiddleware: failed to get sessoin: %s",
						err,
					)
					return
				}

				hu.RespondErrWithStatus(w, http.StatusUnauthorized, Unauthorized)
				return
			}

			if sessionInfo == nil {
				hu.RespondErrWithStatus(w, http.StatusUnauthorized, Unauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), sessionContextKey, sessionInfo)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *middlewareManagerImpl) l() *slog.Logger {
	return mw.logger.WithGroup("middleware-manager")
}
