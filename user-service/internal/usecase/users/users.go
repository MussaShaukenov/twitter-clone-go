package users

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/user-service/internal/dto"
	"MussaShaukenov/twitter-clone-go/user-service/internal/repository"
	"MussaShaukenov/twitter-clone-go/user-service/internal/utils"
	"errors"
	"log"
	"time"
)

type useCase struct {
	userRepo repository.UserRepo
	otpRepo  repository.OTPRepo
}

func NewUserUseCase(userRepo repository.UserRepo, otpRepo repository.OTPRepo) *useCase {
	return &useCase{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (uc *useCase) Register(dto dto.RegisterUserRequest) error {
	if len(dto.FirstName) == 0 {
		log.Println(dto.FirstName)
		return errors.New("first name is empty")
	}
	if len(dto.LastName) == 0 {
		return errors.New("last name is empty")
	}
	if len(dto.Email) == 0 {
		return errors.New("email is empty")
	}
	if dto.Age < 14 {
		return errors.New("come later")
	}

	user := domain.ConvertFromDto(0, dto.FirstName, dto.LastName, dto.Email, dto.Username, dto.Password, dto.Age)

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Call repository to insert user
	return uc.userRepo.Insert(user)
}

func (uc *useCase) Authorize(input dto.LoginRequest) (string, error) {
	if len(input.Username) == 0 {
		return "", errors.New("invalid username")
	}
	// Get user by username
	user, err := uc.userRepo.GetByUsername(input.Username)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			return "", domain.ErrRecordNotFound
		}
		return "", err
	}

	// Validate the password
	if !utils.CheckPassword(input.Password, user.Password) {
		return "", domain.ErrInvalidCredentials
	}

	// Check if it's user's first login
	isFirstLogin, err := uc.userRepo.IsFirstLogin(user.ID)
	if err != nil {
		return "", err
	}

	if isFirstLogin {
		code := uc.Authorize2FA(user.Email)
		return code, nil
	}

	return uc.generateAndStoreSession(user.ID)
}

func (uc *useCase) Authorize2FA(username string) string {
	// Get user email
	user, err := uc.userRepo.GetByUsername(username)

	code := utils.GenerateRandomCode(6)

	// For now, just log OTP
	log.Printf("2FA code for %s: %s\n", user.Email, code)

	// Store the OTP
	err = uc.otpRepo.StoreOTP(user.Email, code)
	if err != nil {
		log.Println("Error storing OTP:", err)
	}
	return code
}

func (uc *useCase) VerifyOTP(email, otp string) (string, error) {
	// retrieve the OTP
	storedOtp, err := uc.otpRepo.GetStoreOTP(email)
	if err != nil {
		return "", errors.New("failed to retrieve OTP")
	}

	// compare the OTPs
	if storedOtp != otp {
		return "", errors.New("invalid OTP")
	}

	// Generate and return session token
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	return uc.generateAndStoreSession(user.ID)
}

func (uc *useCase) generateAndStoreSession(userID int) (string, error) {
	sessionToken, err := utils.GenerateSessionToken(userID)
	if err != nil {
		return "", err
	}
	err = uc.otpRepo.CreateSession(userID, sessionToken, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (uc *useCase) Logout(sessionToken string) error {
	// Invalidate the session token in the repository (if using a session store)
	return uc.otpRepo.DeleteSession(sessionToken)
}

func (uc *useCase) List() ([]*domain.User, error) {
	return uc.userRepo.List()
}
