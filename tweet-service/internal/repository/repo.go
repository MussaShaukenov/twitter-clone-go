package repository

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/domain"
	"context"
)

type TweetRepository interface {
	Insert(in *domain.Tweet) error
	Get(id int64) (*domain.Tweet, error)
	Update(in *domain.Tweet) (*domain.Tweet, error)
	List() ([]*domain.Tweet, error)
	Delete(id int) error
	GetUserTweets(id int) ([]*domain.Tweet, error)
}

type TweetTagRepository interface {
	AddTag(tweetId int64, tagId int64) error
	GetTweetTags(tweetId int64) ([]*domain.Tag, error)
	ListTags() ([]*domain.Tag, error)
}

type TweetStatsRepo interface {
	GetTweetStats(ctx context.Context, tweetID int64) (*domain.TweetStats, error)
	UpdateLikes(ctx context.Context, tweetID int64, likesChange int64) error
	UpdateDislikes(ctx context.Context, tweetID int64, dislikesChange int64) error
}
