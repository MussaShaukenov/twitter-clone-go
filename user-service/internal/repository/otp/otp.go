package otp

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type repo struct {
	Redis *redis.Client
}

func NewOTPRepo(redis *redis.Client) *repo {
	return &repo{
		Redis: redis,
	}
}

func (r *repo) CreateSession(userID int, token string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	value := strconv.Itoa(userID)
	err := r.Redis.Set(context.Background(), key, value, ttl).Err()
	return err
}

func (r *repo) DeleteSession(token string) error {
	key := fmt.Sprintf("session:%s", token)
	err := r.Redis.Del(context.Background(), key).Err()
	return err
}

func (r *repo) StoreOTP(email, code string) error {
	key := fmt.Sprintf("otp:%s", email)
	err := r.Redis.Set(context.Background(), key, code, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store otp: %w", err)
	}
	return nil
}

func (r *repo) GetStoreOTP(email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	otp, err := r.Redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", errors.New("OTP not found or expired")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve OTP: %w", err)
	}
	return otp, nil
}
