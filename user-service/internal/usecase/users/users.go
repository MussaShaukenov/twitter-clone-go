package users

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/user-service/internal/dto"
	"MussaShaukenov/twitter-clone-go/user-service/internal/repository"
	"MussaShaukenov/twitter-clone-go/user-service/internal/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type useCase struct {
	userRepo repository.UserRepo
	otpRepo  repository.OTPRepo
	logger   *zap.SugaredLogger
}

func NewUserUseCase(userRepo repository.UserRepo, otpRepo repository.OTPRepo, logger *zap.SugaredLogger) *useCase {
	return &useCase{
		userRepo: userRepo,
		otpRepo:  otpRepo,
		logger:   logger,
	}
}

func (uc *useCase) Register(dto dto.RegisterUserRequest) error {
	// Validate input
	err := validateDtoInput(dto)
	if err != nil {
		return err
	}

	user := domain.ConvertFromDto(0, dto.FirstName, dto.LastName, dto.Email, dto.Username, dto.Password, dto.Age)

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		uc.logger.Errorw("failed to hash password", "error", err)
		return err
	}
	user.Password = hashedPassword

	// Insert user into repository
	uc.logger.Info("Inserting user into repository")
	if err = uc.userRepo.Insert(user); err != nil {
		uc.logger.Errorw("Failed to insert user", "error", err)
		return err
	}
	return nil
}

func (uc *useCase) Authorize(input dto.LoginRequest) (string, error) {
	// Validate input
	if err := validateLoginInput(input); err != nil {
		uc.logger.Errorw("Validation failed for Login request", "error", err)
		return "", err
	}
	// Get user by username
	user, err := uc.userRepo.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			uc.logger.Warnw("User not found during authorization", "username", input.Username)
			return "", domain.ErrRecordNotFound
		default:
			return "", err
		}
	}

	// Validate the password
	if !utils.CheckPassword(input.Password, user.Password) {
		uc.logger.Warn("Invalid credentials provided")
		return "", domain.ErrInvalidCredentials
	}

	// Check if it's user's first login
	isFirstLogin, err := uc.userRepo.IsFirstLogin(user.ID)
	if err != nil {
		uc.logger.Errorw("Failed to check first login", "error", err)
		return "", err
	}

	if isFirstLogin {
		code := uc.Authorize2FA(user.Email)
		return code, nil
	}

	token, err := uc.generateAndStoreSession(user.ID)
	if err != nil {
		uc.logger.Errorw("Failed to generate and store session", "error", err)
		return "", err
	}
	return token, nil
}

func (uc *useCase) Authorize2FA(username string) string {
	// Get user email
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil {
		uc.logger.Errorw("Failed to fetch user by username for 2FA", "username", username, "error", err)
		return ""
	}

	code := utils.GenerateRandomCode(6)
	uc.logger.Infow("Generated 2FA code", "email", user.Email, "code", code)

	// Store the OTP
	err = uc.otpRepo.StoreOTP(user.Email, code)
	if err != nil {
		uc.logger.Errorw("Failed to store OTP", "email", user.Email, "error", err)
	}
	return code
}

func (uc *useCase) VerifyOTP(email, otp string) (string, error) {
	// retrieve the OTP
	storedOtp, err := uc.otpRepo.GetStoreOTP(email)
	if err != nil {
		uc.logger.Errorw("Failed to retrieve OTP", "email", email, "error", err)
		return "", ErrFailedToRetrieveOTP
	}

	// compare the OTPs
	if storedOtp != otp {
		uc.logger.Warnw("Invalid OTP provided", "email", email)
		return "", ErrInvalidOTP
	}

	// Generate and return session token
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		uc.logger.Errorw("Failed to fetch user by email for session generation", "email", email, "error", err)
		return "", err
	}
	token, err := uc.generateAndStoreSession(user.ID)
	if err != nil {
		uc.logger.Errorw("Failed to generate and store session", "error", err)
		return "", err
	}
	return token, nil
}

func (uc *useCase) generateAndStoreSession(userID int) (string, error) {
	sessionToken, err := utils.GenerateSessionToken(userID)
	if err != nil {
		uc.logger.Errorw("Failed to generate session token", "userID", userID, "error", err)
		return "", err
	}
	err = uc.otpRepo.CreateSession(userID, sessionToken, 24*time.Hour)
	if err != nil {
		uc.logger.Errorw("Failed to create session", "userID", userID, "error", err)
		return "", err
	}
	return sessionToken, nil
}

func (uc *useCase) Logout(sessionToken string) error {
	// Invalidate the session token in the repository (if using a session store)
	if err := uc.otpRepo.DeleteSession(sessionToken); err != nil {
		uc.logger.Errorw("Failed to delete session", "sessionToken", sessionToken, "error", err)
		return err
	}
	return nil
}

func (uc *useCase) List() ([]*domain.User, error) {
	users, err := uc.userRepo.List()
	if err != nil {
		uc.logger.Errorw("Failed to fetch user list", "error", err)
		return nil, err
	}
	return users, nil
}

func validateDtoInput(dto dto.RegisterUserRequest) error {
	if err := validateNonEmptyField("First Name", dto.FirstName); err != nil {
		return err
	}
	if err := validateNonEmptyField("Last Name", dto.LastName); err != nil {
		return err
	}
	if err := validateNonEmptyField("Email", dto.Email); err != nil {
		return err
	}
	if dto.Age < 14 {
		return ErrAgeRestrict
	}
	return nil
}

func validateLoginInput(input dto.LoginRequest) error {
	return validateNonEmptyField("Username", input.Username)
}

func validateNonEmptyField(fieldName, value string) error {
	if len(value) == 0 {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}
