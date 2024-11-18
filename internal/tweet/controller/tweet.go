package controller

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet/dto"
	"MussaShaukenov/twitter-clone-go/internal/tweet/usecase"
	"MussaShaukenov/twitter-clone-go/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type controller struct {
	service usecase.TweetUseCase
}

func NewController(service usecase.TweetUseCase) *controller {
	return &controller{
		service: service,
	}
}

func (c *controller) CreateTweetHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.TweetDto

	// Read Body parameters
	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call useCase
	err = c.service.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]string{"message": "tweet created successfully"}
	err = utils.WriteJson(w, http.StatusCreated, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *controller) GetTweetByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Read ID
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// GetTweet record
	tweet, err := c.service.Get(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	err = utils.WriteJson(w, http.StatusOK, tweet, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (c *controller) ListTweetsHandler(w http.ResponseWriter, r *http.Request) {
	// GetTweet records
	tweets, err := c.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, tweets, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *controller) UpdateTweetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// GetTweet record
	tweet, err := c.service.Get(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var in dto.UpdateTweetDto
	err = utils.ReadJson(w, r, &in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update existing tweet fields if input is not nil
	if in.Title != nil {
		tweet.Title = *in.Title
	}
	if in.Content != nil {
		tweet.Content = *in.Content
	}
	if in.Topic != nil {
		tweet.Topic = *in.Topic
	}

	// UpdateTweets record
	res, err := c.service.Update(*tweet)
	if err != nil {
		http.Error(w, "error 1", http.StatusInternalServerError)
		return
	}

	// Return response
	err = utils.WriteJson(w, http.StatusOK, res, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (c *controller) DeleteTweetHandler(w http.ResponseWriter, r *http.Request) {
	// GetTweet ID
	id, err := utils.GetIdFromQueryParam(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	_, err = c.service.Get(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// DeleteTweet record
	err = c.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	err = utils.WriteJson(w, http.StatusOK, `{"message": "tweet successfully deleted"}`, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *controller) GetUserTweets(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tweets, err := c.service.GetUserTweets(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, tweets, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
