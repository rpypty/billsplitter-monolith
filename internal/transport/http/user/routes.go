package user

import (
	"billsplitter-monolith/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r chi.Router, ctrl Controller, mw middleware.Manager) {
	// User routes with user ID in path
	r.Route("/user", func(r chi.Router) {
		// Routes requiring authentication
		r.With(mw.Auth()).Group(func(r chi.Router) {
			// User profile routes
			r.Put("/{id}/profile", ctrl.UpdateUserProfile)

			// Payment methods routes
			r.Get("/{id}/payment_methods", ctrl.GetPaymentMethods)                 // GET /user/{id}/payment_methods
			r.Post("/{id}/payment_methods", ctrl.CreatePaymentMethod)              // POST /user/{id}/payment_methods
			r.Put("/{id}/payment_methods/{methodId}", ctrl.UpdatePaymentMethod)    // PUT /user/{id}/payment_methods/{methodId}
			r.Delete("/{id}/payment_methods/{methodId}", ctrl.DeletePaymentMethod) // DELETE /user/{id}/payment_methods/{methodId}
		})
	})

	// Payment methods routes without user ID in path (for frontend compatibility)
	r.Route("/payment_methods", func(r chi.Router) {
		// Routes requiring authentication
		r.With(mw.Auth()).Group(func(r chi.Router) {
			r.Get("/", ctrl.GetPaymentMethods)                 // GET /payment_methods
			r.Post("/", ctrl.CreatePaymentMethod)              // POST /payment_methods
			r.Put("/{methodId}", ctrl.UpdatePaymentMethod)     // PUT /payment_methods/{methodId}
			r.Delete("/{methodId}", ctrl.DeletePaymentMethod)  // DELETE /payment_methods/{methodId}
		})
	})
}
