package user

import (
	"billsplitter-monolith/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, ctrl Controller, mw middleware.Manager) {
	r.Route("/user", func(r chi.Router) {
		// public routes
		r.With(mw.Auth()).Put("/{id}/profile", ctrl.UpdateUserProfile)
	})
}
