package controller

import (
	followerCtrl "MussaShaukenov/twitter-clone-go/user-service/internal/controller/followers"
	userCtrl "MussaShaukenov/twitter-clone-go/user-service/internal/controller/users"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(ctrl *userCtrl.UserController) http.Handler {
	router := chi.NewRouter()

	router.Post("/register", ctrl.RegisterHandler)
	router.Post("/authorize", ctrl.AuthorizeHandler)
	router.Post("/authorize2fa", ctrl.Authorize2FAHandler)
	router.Post("/verifyotp", ctrl.VerifyOTPHandler)
	router.Post("/logout", ctrl.LogoutHandler)
	router.Get("/", ctrl.ListHandler)

	return router
}

func RegisterFollowerRoutes(ctrl *followerCtrl.FollowerController) http.Handler {
	router := chi.NewRouter()

	router.Post("/follow", ctrl.FollowHandler)
	router.Post("/unfollow", ctrl.UnfollowHandler)
	router.Get("/followers/{id}", ctrl.GetFollowersHandler)
	router.Get("/following/{id}", ctrl.GetFollowingHandler)
	router.Get("/isfollowing", ctrl.IsFollowingHandler)

	return router
}
