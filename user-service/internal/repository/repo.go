package repository

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"time"
)

type UserRepo interface {
	Insert(in *domain.User) error
	Delete(id int) error
	GetByUsername(username string) (*domain.User, error)
	GetUserEmail(id int) (string, error)
	GetByEmail(email string) (*domain.User, error)
	IsFirstLogin(userId int) (bool, error)
	List() ([]*domain.User, error)
}

type OTPRepo interface {
	CreateSession(userID int, token string, ttl time.Duration) error
	DeleteSession(token string) error
	StoreOTP(email, code string) error
	GetStoreOTP(email string) (string, error)
}
