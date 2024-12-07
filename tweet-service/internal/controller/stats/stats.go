package stats

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/usecase"
	"MussaShaukenov/twitter-clone-go/tweet-service/pkg/utils"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type TweetStatsController struct {
	useCase usecase.TweetStatsUseCase
}

func NewTweetStatsController(useCase usecase.TweetStatsUseCase) *TweetStatsController {
	return &TweetStatsController{
		useCase: useCase,
	}
}

func (c *TweetStatsController) GetTweetStatsHandler(w http.ResponseWriter, r *http.Request) {
	tweetID, err := getTweetID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := c.useCase.GetTweetStats(context.Background(), tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, stats, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *TweetStatsController) AddLikeHandler(w http.ResponseWriter, r *http.Request) {
	tweetID, err := getTweetID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCase.AddLike(context.Background(), tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *TweetStatsController) AddDislikeHandler(w http.ResponseWriter, r *http.Request) {
	tweetID, err := getTweetID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCase.AddDislike(context.Background(), tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *TweetStatsController) RemoveLikeHandler(w http.ResponseWriter, r *http.Request) {
	tweetID, err := getTweetID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCase.RemoveLike(context.Background(), tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *TweetStatsController) RemoveDislikeHandler(w http.ResponseWriter, r *http.Request) {
	tweetID, err := getTweetID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCase.RemoveDislike(context.Background(), tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTweetID(r *http.Request) (int64, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
