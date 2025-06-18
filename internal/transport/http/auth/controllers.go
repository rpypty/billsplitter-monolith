package auth

import (
	"net/http"

	"billsplitter-monolith/internal/domain/auth"
)

type Controller struct {
	svc auth.Service
}

func NewController(svc auth.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (c *Controller) LoginTelegram(w http.ResponseWriter, r *http.Request) {

}
