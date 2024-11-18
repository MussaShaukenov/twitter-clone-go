package domain

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
