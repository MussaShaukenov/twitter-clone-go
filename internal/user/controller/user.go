package controller

import (
	"MussaShaukenov/twitter-clone-go/internal/user/dto"
	"MussaShaukenov/twitter-clone-go/internal/user/usecase"
	"MussaShaukenov/twitter-clone-go/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
)

type UserController struct {
	useCase usecase.UserUseCase
}

func NewController(uc usecase.UserUseCase) *UserController {
	return &UserController{
		useCase: uc,
	}
}

func (ctrl *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterUserRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("incoming input:", input)

	err = ctrl.useCase.Register(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]string{"message": "successfully registered"}
	err = utils.WriteJson(w, http.StatusCreated, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *UserController) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest

	// Read body params
	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ctrl.useCase.Authorize(input)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"message": "successfully authorized",
		"token":   token,
	}

	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *UserController) Authorize2FAHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := ctrl.useCase.Authorize2FA(input.Username)
	response := map[string]string{
		"message": "successfully authorized",
		"token":   token,
	}

	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *UserController) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.VerifyOTPRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ctrl.useCase.VerifyOTP(input.Email, input.OTP)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"message": "successfully authorized",
		"token":   token,
	}

	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (ctrl *UserController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	err := ctrl.useCase.Logout(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "logout successful",
	})
}
