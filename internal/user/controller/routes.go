package controller

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterUserRoutes(ctrl *UserController) http.Handler {
	router := chi.NewRouter()

	router.Post("/register", ctrl.RegisterHandler)
	router.Post("/authorize", ctrl.AuthorizeHandler)
	router.Post("/logout", ctrl.LogoutHandler)

	return router
}
