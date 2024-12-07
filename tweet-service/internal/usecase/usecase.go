package usecase

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/domain"
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/dto"
	"context"
)

type TweetUseCase interface {
	Create(dto dto.TweetDto) error
	Get(id int64) (*dto.TweetDto, error)
	List() ([]*dto.TweetDto, error)
	Update(in dto.TweetDto) (*dto.GetTweetResponse, error)
	Delete(id int) error
	GetUserTweets(id int) ([]*dto.TweetDto, error)
}

type TweetStatsUseCase interface {
	GetTweetStats(ctx context.Context, tweetID int64) (*domain.TweetStats, error)
	AddLike(ctx context.Context, tweetID int64) error
	AddDislike(ctx context.Context, tweetID int64) error
	RemoveLike(ctx context.Context, tweetID int64) error
	RemoveDislike(ctx context.Context, tweetID int64) error
}

type TweetTagUseCase interface {
	AddTag(tweetId int64, tagId int64) error
	GetTweetTags(tweetId int64) ([]*dto.TagDto, error)
	ListTags() ([]*dto.TagDto, error)
}
