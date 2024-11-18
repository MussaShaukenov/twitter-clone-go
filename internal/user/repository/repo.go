package repository

import "MussaShaukenov/twitter-clone-go/internal/user/domain"

type UserRepo interface {
	Insert(in *domain.User) error
	Get(id int) (*domain.User, error)
	Delete(id int) error
	GetByUsername(username string) (*domain.User, error)
	CreateSession(userID int, token string) error
	DeleteSession(token string) error
}
