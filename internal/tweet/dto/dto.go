package dto

import "time"

type TweetDto struct {
	ID        int       `json:"id"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Topic     string    `json:"topic,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty,omitempty"`
}

type UpdateTweetDto struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
	Topic   *string `json:"topic,omitempty"`
}

type GetTweetResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Topic     string    `json:"topic"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
