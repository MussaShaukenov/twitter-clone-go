package users

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/dto"
	"MussaShaukenov/twitter-clone-go/user-service/internal/usecase"
	"MussaShaukenov/twitter-clone-go/user-service/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	useCase usecase.UserUseCase
	logger  *zap.SugaredLogger
}

func NewUserController(userUC usecase.UserUseCase, logger *zap.SugaredLogger) *UserController {
	return &UserController{
		useCase: userUC,
		logger:  logger,
	}
}

func (ctrl *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterUserRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		ctrl.logger.Errorw("failed to read json", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctrl.logger.Infow("incoming input", "input", input)

	err = ctrl.useCase.Register(input)
	if err != nil {
		ctrl.logger.Errorw("failed to register user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]string{"message": "successfully registered"}
	err = utils.WriteJson(w, http.StatusCreated, response, nil)
	if err != nil {
		ctrl.logger.Errorw("failed to write json", "error", err)
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

	ctrl.logger.Infow("incoming input", "input", input)

	token, err := ctrl.useCase.Authorize(input)
	if err != nil {
		ctrl.logger.Errorw("failed to authorize", "error", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"message": "successfully authorized",
		"token":   token,
	}

	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		ctrl.logger.Errorw("failed to write json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *UserController) Authorize2FAHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		ctrl.logger.Errorw("failed to read json", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctrl.logger.Infow("incoming input", "input", input)

	otp := ctrl.useCase.Authorize2FA(input.Username)
	response := map[string]string{
		"message": "successfully authorized",
		"otp":     otp,
	}

	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		ctrl.logger.Errorw("failed to write json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *UserController) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.VerifyOTPRequest

	err := utils.ReadJson(w, r, &input)
	if err != nil {
		ctrl.logger.Errorw("failed to read json", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctrl.logger.Infow("incoming input", "input", input)

	token, err := ctrl.useCase.VerifyOTP(input.Email, input.OTP)
	if err != nil {
		ctrl.logger.Errorw("failed to verify otp", "error", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"message": "successfully authorized",
		"token":   token,
	}

	err = utils.WriteJson(w, http.StatusOK, response, nil)
	if err != nil {
		ctrl.logger.Errorw("failed to write json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (ctrl *UserController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		ctrl.logger.Error("token is required")
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	err := ctrl.useCase.Logout(token)
	if err != nil {
		ctrl.logger.Errorw("failed to logout", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, map[string]interface{}{"message": "logout successful"}, nil)
	if err != nil {
		ctrl.logger.Errorw("failed to write json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ctrl *UserController) ListHandler(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.useCase.List()
	if err != nil {
		ctrl.logger.Errorw("failed to list users", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.WriteJson(w, http.StatusOK, users, nil)
	if err != nil {
		ctrl.logger.Errorw("failed to write json", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
