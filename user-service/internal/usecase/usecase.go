package usecase

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/user-service/internal/dto"
)

type UserUseCase interface {
	Register(dto dto.RegisterUserRequest) error
	Authorize(input dto.LoginRequest) (string, error)
	Logout(sessionToken string) error
	Authorize2FA(email string) string
	VerifyOTP(email, otp string) (string, error)
	List() ([]*domain.User, error)
}

type FollowerUseCase interface {
	Follow(followerID, followeeID int) error
	Unfollow(followerID, followeeID int) error
	GetFollowers(userID int) ([]*domain.User, error)
	GetFollowing(userID int) ([]*domain.User, error)
	IsFollowing(followerID, followeeID int) (bool, error)
}
