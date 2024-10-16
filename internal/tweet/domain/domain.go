package domain

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet/dto"
	"time"
)

type Tweet struct {
	ID        int64
	Title     string
	Topic     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ConvertFromDto(id int, title, content, topic string) *Tweet {
	return &Tweet{
		ID:        int64(id),
		Title:     title,
		Content:   content,
		Topic:     topic,
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
