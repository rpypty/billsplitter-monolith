package meet

import (
	"billsplitter-monolith/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, ctrl Controller, mw middleware.Manager) {
	r.Route("/meet", func(r chi.Router) {
		// public routes
		r.With(mw.Auth()).Post("/", ctrl.CreateMeet)
	})
}
