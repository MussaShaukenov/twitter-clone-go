package domain

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/dto"
	"time"
)

type Tweet struct {
	ID        int64
	Title     string
	Topic     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    int
	Tags      []Tag
}

type Tag struct {
	ID   int64
	Name string
}

type TweetStats struct {
	TweetID    int64     `bson:"tweet_id"`
	Likes      int64     `bson:"likes"`
	Dislikes   int64     `bson:"dislikes"`
	LastUpdate time.Time `bson:"last_update"`
}

func ConvertFromDto(id int, title, content, topic string, userId int) *Tweet {
	return &Tweet{
		ID:        int64(id),
		Title:     title,
		Content:   content,
		Topic:     topic,
		UserId:    userId,
		CreatedAt: time.Now(),
	}
}

func ConvertToDto(tweet *Tweet) *dto.TweetDto {
	return &dto.TweetDto{
		Title:   tweet.Title,
		Content: tweet.Content,
		Topic:   tweet.Topic,
	}
}

func ConvertToGetTweetResponseDto(tweet *Tweet) *dto.GetTweetResponse {
	return &dto.GetTweetResponse{
		ID:        tweet.ID,
		Title:     tweet.Title,
		Content:   tweet.Content,
		Topic:     tweet.Topic,
		CreatedAt: tweet.CreatedAt,
		UpdatedAt: tweet.UpdatedAt,
	}
}
