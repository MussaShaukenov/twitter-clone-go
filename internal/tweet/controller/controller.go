package controller

import "net/http"

type TweetController interface {
	CreateTweetHandler(w http.ResponseWriter, r *http.Request)
	GetTweetByIdHandler(w http.ResponseWriter, r *http.Request)
	ListTweetsHandler(w http.ResponseWriter, r *http.Request)
	UpdateTweetHandler(w http.ResponseWriter, r *http.Request)
	DeleteTweetHandler(w http.ResponseWriter, r *http.Request)
	GetUserTweets(w http.ResponseWriter, r *http.Request)
}
