package domain

import (
	"time"
)

type User struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	Age          int
	Username     string
	Password     string
	IsFirstLogin bool // New field to track first login
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
