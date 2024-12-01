package controller

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterTweetRoutes(ctrl TweetController) http.Handler {
	router := chi.NewRouter()

	router.Post("/", ctrl.CreateTweetHandler)
	router.Get("/{id}", ctrl.GetTweetByIdHandler)
	router.Get("/", ctrl.ListTweetsHandler)
	router.Patch("/{id}", ctrl.UpdateTweetHandler)
	router.Delete("/{id}", ctrl.DeleteTweetHandler)
	router.Get("/{user_id}", ctrl.GetUserTweetsHandler)

	router.Post("/{tweet_id}/tags", ctrl.AddTweetTagHandler)
	router.Get("/{tweet_id}/tags", ctrl.GetTweetTagsHandler)

	return router
}

func RegisterTagsRoutes(ctrl TweetTagController) http.Handler {
	router := chi.NewRouter()
	router.Get("/", ctrl.ListTagsHandler)
	return router
}
