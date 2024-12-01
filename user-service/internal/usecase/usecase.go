package usecase

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/user-service/internal/dto"
	"MussaShaukenov/twitter-clone-go/user-service/internal/repository"
	"MussaShaukenov/twitter-clone-go/user-service/internal/utils"
	"errors"
	"log"
	"time"
)

type UserUseCase interface {
	Register(dto dto.RegisterUserRequest) error
	Authorize(input dto.LoginRequest) (string, error)
	Logout(sessionToken string) error
	Authorize2FA(email string) string
	VerifyOTP(email, otp string) (string, error)
	List() ([]*domain.User, error)
}

type userUseCase struct {
	repo      repository.UserRepo
	redisRepo repository.OTPRepo
}

func NewUserUseCase(repo repository.UserRepo, redisRepo repository.OTPRepo) *userUseCase {
	return &userUseCase{
		repo:      repo,
		redisRepo: redisRepo,
	}
}

func (uc *userUseCase) Register(dto dto.RegisterUserRequest) error {
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
	return uc.repo.Insert(user)
}

func (uc *userUseCase) Authorize(input dto.LoginRequest) (string, error) {
	if len(input.Username) == 0 {
		return "", errors.New("invalid username")
	}
	// Get user by username
	user, err := uc.repo.GetByUsername(input.Username)
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
	isFirstLogin, err := uc.repo.IsFirstLogin(user.ID)
	if err != nil {
		return "", err
	}

	if isFirstLogin {
		code := uc.Authorize2FA(user.Email)
		return code, nil
	}

	return uc.generateAndStoreSession(user.ID)
}

func (uc *userUseCase) Authorize2FA(username string) string {
	// Get user email
	user, err := uc.repo.GetByUsername(username)

	code := utils.GenerateRandomCode(6)

	// For now, just log OTP
	log.Printf("2FA code for %s: %s\n", user.Email, code)

	// Store the OTP
	err = uc.redisRepo.StoreOTP(user.Email, code)
	if err != nil {
		log.Println("Error storing OTP:", err)
	}
	return code
}

func (uc *userUseCase) VerifyOTP(email, otp string) (string, error) {
	// retrieve the OTP
	storedOtp, err := uc.redisRepo.GetStoreOTP(email)
	if err != nil {
		return "", errors.New("failed to retrieve OTP")
	}

	// compare the OTPs
	if storedOtp != otp {
		return "", errors.New("invalid OTP")
	}

	// Generate and return session token
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	return uc.generateAndStoreSession(user.ID)
}

func (uc *userUseCase) generateAndStoreSession(userID int) (string, error) {
	sessionToken, err := utils.GenerateSessionToken(userID)
	if err != nil {
		return "", err
	}
	err = uc.redisRepo.CreateSession(userID, sessionToken, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (uc *userUseCase) Logout(sessionToken string) error {
	// Invalidate the session token in the repository (if using a session store)
	return uc.redisRepo.DeleteSession(sessionToken)
}

func (uc *userUseCase) List() ([]*domain.User, error) {
	return uc.repo.List()
}
