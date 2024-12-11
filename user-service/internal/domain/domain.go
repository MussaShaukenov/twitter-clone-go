package domain

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	Age          int       `json:"age"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	IsFirstLogin bool      `json:"isFirstLogin"` // New field to track first login
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Follower struct {
	FollowerID int // ID of the user who is following
	FollowedID int // ID of the user who is being followed
	CreatedAt  time.Time
}

func ConvertFromDto(id int, firstName, lastName, email, username, password string, age int) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Age:       age,
		Username:  username,
		Password:  password,
	}
}
