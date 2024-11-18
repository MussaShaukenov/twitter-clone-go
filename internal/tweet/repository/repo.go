package repository

import "MussaShaukenov/twitter-clone-go/internal/tweet/domain"

type TweetRepository interface {
	Insert(in *domain.Tweet) error
	Get(id int64) (*domain.Tweet, error)
	Update(in *domain.Tweet) (*domain.Tweet, error)
	List() ([]*domain.Tweet, error)
	Delete(id int) error
	GetUserTweets(id int) ([]*domain.Tweet, error)
}
