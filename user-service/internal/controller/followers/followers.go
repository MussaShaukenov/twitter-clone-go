package followers

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/dto"
	"MussaShaukenov/twitter-clone-go/user-service/internal/usecase"
	"MussaShaukenov/twitter-clone-go/user-service/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type FollowerController struct {
	useCase usecase.FollowerUseCase
}

func NewFollowerController(followerUC usecase.FollowerUseCase) *FollowerController {
	return &FollowerController{
		useCase: followerUC,
	}
}

func (ctrl *FollowerController) FollowHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.FollowRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctrl.useCase.Follow(input.FollowerID, input.FollowedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "successfully followed"}
	err = utils.WriteJson(w, http.StatusCreated, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *FollowerController) UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.FollowRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctrl.useCase.Unfollow(input.FollowerID, input.FollowedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "successfully unfollowed"}
	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *FollowerController) GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users, err := ctrl.useCase.GetFollowers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, users, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *FollowerController) GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users, err := ctrl.useCase.GetFollowing(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, users, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *FollowerController) IsFollowingHandler(w http.ResponseWriter, r *http.Request) {
	followerID, err := strconv.Atoi(chi.URLParam(r, "followerID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	followedID, err := strconv.Atoi(chi.URLParam(r, "followedID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isFollowing, err := ctrl.useCase.IsFollowing(followerID, followedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, map[string]bool{"is_following": isFollowing}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
