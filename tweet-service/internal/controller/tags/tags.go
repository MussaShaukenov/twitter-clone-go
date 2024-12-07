package tags

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/usecase"
	"MussaShaukenov/twitter-clone-go/tweet-service/pkg/utils"
	"net/http"
)

type TweetTags struct {
	useCase usecase.TweetTagUseCase
}

func NewTweetTagsController(useCase usecase.TweetTagUseCase) *TweetTags {
	return &TweetTags{
		useCase: useCase,
	}
}

func (c *TweetTags) AddTweetTagHandler(w http.ResponseWriter, r *http.Request) {
	tweetId, err := utils.GetIdFromQueryParam(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tagId, err := utils.GetIdFromQueryParam(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCase.AddTag(int64(tweetId), int64(tagId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, `{"message": "tag added to tweet"}`, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *TweetTags) GetTweetTagsHandler(w http.ResponseWriter, r *http.Request) {
	tweetId, err := utils.GetIdFromQueryParam(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tags, err := c.useCase.GetTweetTags(int64(tweetId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, tags, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *TweetTags) ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := c.useCase.ListTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, tags, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
