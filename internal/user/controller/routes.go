package controller

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterUserRoutes(ctrl *UserController) http.Handler {
	router := chi.NewRouter()

	router.Post("/register", ctrl.RegisterHandler)
	router.Post("/authorize", ctrl.AuthorizeHandler)
	router.Post("/authorize2fa", ctrl.Authorize2FAHandler)
	router.Post("/verifyotp", ctrl.VerifyOTPHandler)
	router.Post("/logout", ctrl.LogoutHandler)

	return router
}
