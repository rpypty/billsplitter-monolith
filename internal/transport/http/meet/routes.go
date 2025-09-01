package meet

import (
	"billsplitter-monolith/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, ctrl Controller, mw middleware.Manager) {
	r.Route("/meets", func(r chi.Router) {
		// public routes
		r.With(mw.Auth()).Post("/", ctrl.CreateMeet)

		r.With(mw.Auth()).Get("/", ctrl.FetchUserMeets)

		r.With(mw.Auth()).Get("/{id}", ctrl.GetMeetDetailsByID)
	})
}
