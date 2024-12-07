package otp

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type repository struct {
	redis  *redis.Client
	logger *zap.SugaredLogger
}

func NewOTPRepo(redis *redis.Client, logger *zap.SugaredLogger) *repository {
	return &repository{
		redis:  redis,
		logger: logger,
	}
}

func (repo *repository) CreateSession(userID int, token string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	value := strconv.Itoa(userID)
	err := repo.redis.Set(context.Background(), key, value, ttl).Err()
	if err != nil {
		repo.logger.Errorw("failed to create session", "error", err)
		return ErrFailedToCreateSession
	}
	return nil
}

func (repo *repository) DeleteSession(token string) error {
	key := fmt.Sprintf("session:%s", token)
	err := repo.redis.Del(context.Background(), key).Err()
	if err != nil {
		repo.logger.Errorw("failed to delete session", "error", err)
		return fmt.Errorf("failed to delete session: %w", err)

	}
	return nil
}

func (repo *repository) StoreOTP(email, code string) error {
	key := fmt.Sprintf("otp:%s", email)
	err := repo.redis.Set(context.Background(), key, code, 5*time.Minute).Err()
	if err != nil {
		repo.logger.Errorw("failed to store otp", "error", err)
		return ErrFailedToStoreOTP
	}
	return nil
}

func (repo *repository) GetStoreOTP(email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	otp, err := repo.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		repo.logger.Warnw("OTP not found or expired", "email", email)
		return "", ErrOTPNotFound
	} else if err != nil {
		return "", ErrFailedToRetrieveOTP
	}
	return otp, nil
}
