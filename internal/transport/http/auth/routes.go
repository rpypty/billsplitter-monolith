package auth

import (
	"billsplitter-monolith/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, ctrl Controller, mw middleware.Manager) {
	r.Route("/auth", func(r chi.Router) {
		// public routes
		r.Post("/login/telegram", ctrl.LoginTelegram)

		// auth routes
		r.With(mw.Auth()).Get("/me", ctrl.Me)
	})
}
