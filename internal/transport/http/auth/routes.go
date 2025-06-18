package auth

import "github.com/go-chi/chi/v5"

func InitRoutes(r chi.Router, ctrl *Controller) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/login/telegram", ctrl.LoginTelegram)
	})
}
