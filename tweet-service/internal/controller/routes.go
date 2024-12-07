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

	return router
}

func RegisterTagsRoutes(ctrl TweetTagController) http.Handler {
	router := chi.NewRouter()

	router.Get("/", ctrl.ListTagsHandler)
	router.Post("/{tweet_id}/tags", ctrl.AddTweetTagHandler)
	router.Get("/{tweet_id}/tags", ctrl.GetTweetTagsHandler)

	return router
}

func RegisterStatsRoutes(ctrl TweetStatsController) http.Handler {
	router := chi.NewRouter()

	router.Get("/{tweet_id}/stats", ctrl.GetTweetStatsHandler)
	router.Post("/{tweet_id}/like", ctrl.AddLikeHandler)
	router.Post("/{tweet_id}/dislike", ctrl.AddDislikeHandler)
	router.Delete("/{tweet_id}/like", ctrl.RemoveLikeHandler)
	router.Delete("/{tweet_id}/dislike", ctrl.RemoveDislikeHandler)

	return router
}
