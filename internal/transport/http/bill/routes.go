package bill

import (
	"billsplitter-monolith/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, ctrl Controller, mw middleware.Manager) {
	r.Route("/bills", func(r chi.Router) {
		r.With(mw.Auth()).Post("/", ctrl.CreateBill)
		r.With(mw.Auth()).Get("/", ctrl.FetchEventBills)
	})
}
