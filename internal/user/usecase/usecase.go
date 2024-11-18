package usecase

import (
	"MussaShaukenov/twitter-clone-go/internal/user/domain"
	"MussaShaukenov/twitter-clone-go/internal/user/dto"
	"MussaShaukenov/twitter-clone-go/internal/user/repository"
	"MussaShaukenov/twitter-clone-go/internal/user/utils"
	"errors"
	"log"
)

type UserUseCase interface {
	Register(dto dto.RegisterUserRequest) error
	Authorize(input dto.LoginRequest) (string, error)
	Logout(sessionToken string) error
}

type userUseCase struct {
	repo repository.UserRepo
}

func NewUserUseCase(repo repository.UserRepo) *userUseCase {
	return &userUseCase{
		repo: repo,
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
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	// Validate the password
	if !utils.CheckPassword(input.Password, user.Password) {
		return "", domain.ErrInvalidCredentials
	}

	// Generate a session token
	sessionToken, err := utils.GenerateSessionToken(user.ID)
	if err != nil {
		return "", err
	}

	// Store session in repository (optional, if using a session store)
	err = uc.repo.CreateSession(user.ID, sessionToken)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (uc *userUseCase) Logout(sessionToken string) error {
	// Invalidate the session token in the repository (if using a session store)
	return uc.repo.DeleteSession(sessionToken)
}
