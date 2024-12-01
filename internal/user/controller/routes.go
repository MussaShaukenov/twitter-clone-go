package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(ctrl *UserController) http.Handler {
	router := chi.NewRouter()

	router.Post("/register", ctrl.RegisterHandler)
	router.Post("/authorize", ctrl.AuthorizeHandler)
	router.Post("/authorize2fa", ctrl.Authorize2FAHandler)
	router.Post("/verifyotp", ctrl.VerifyOTPHandler)
	router.Post("/logout", ctrl.LogoutHandler)
	router.Get("/", ctrl.ListHandler)

	return router
}
