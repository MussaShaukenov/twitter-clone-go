package controller

import "net/http"

type TweetController interface {
	CreateTweetHandler(w http.ResponseWriter, r *http.Request)
	GetTweetByIdHandler(w http.ResponseWriter, r *http.Request)
	ListTweetsHandler(w http.ResponseWriter, r *http.Request)
	UpdateTweetHandler(w http.ResponseWriter, r *http.Request)
	DeleteTweetHandler(w http.ResponseWriter, r *http.Request)
	GetUserTweetsHandler(w http.ResponseWriter, r *http.Request)
	TweetTagController
}

type TweetTagController interface {
	AddTweetTagHandler(w http.ResponseWriter, r *http.Request)
	GetTweetTagsHandler(w http.ResponseWriter, r *http.Request)
	ListTagsHandler(w http.ResponseWriter, r *http.Request)
}
