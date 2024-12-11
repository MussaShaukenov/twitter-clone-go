package tweets

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	Db          *pgxpool.Pool
	RedisClient *redis.Client
	CacheTTL    time.Duration
}

func NewTweetRepository(db *pgxpool.Pool, redisClient *redis.Client, cacheTTL time.Duration) *repository {
	return &repository{
		Db:          db,
		RedisClient: redisClient,
		CacheTTL:    cacheTTL,
	}
}

func (pg *repository) RebuildCache() error {
	// fetch all tweets
	tweets, err := pg.List()
	if err != nil {
		return err
	}

	// serialize tweets and store them in Redis
	cachedData, err := json.Marshal(tweets)
	if err != nil {
		return err
	}
	return pg.RedisClient.Set(context.Background(), "tweets:list", cachedData, pg.CacheTTL).Err()
}

func (pg *repository) Insert(in *domain.Tweet) error {
	log.Println("in: ", in)

	query := `
			INSERT INTO tweets (title, content, topic, user_id, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, NOW(), NOW())`

	args := []interface{}{in.Title, in.Content, in.Topic, in.UserId}
	_, err := pg.Db.Exec(context.Background(), query, args...)
	if err != nil {
		return err
	}

	// Delete cache
	if err := pg.RedisClient.Del(context.Background(), "tweets:list").Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	// Invalidate cache
	if err := pg.RebuildCache(); err != nil {
		return fmt.Errorf("failed to rebuild cache: %w", err)
	}

	return nil
}

func (pg *repository) Get(id int64) (*domain.Tweet, error) {
	query := `
				SELECT id, title, content, topic, created_at FROM tweets
				WHERE id = $1`

	var tweet domain.Tweet

	err := pg.Db.QueryRow(context.Background(), query, id).Scan(
		&tweet.ID,
		&tweet.Title,
		&tweet.Content,
		&tweet.Topic,
		&tweet.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, domain.ErrRecordNotFoundX
		default:
			return nil, err
		}
	}

	return &tweet, nil
}

func (pg *repository) List() ([]*domain.Tweet, error) {
	ctx := context.Background()
	cacheKey := "tweets:list"

	// Check if the data is in the cache
	cachedTweets, err := pg.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var tweets []*domain.Tweet
		err := json.Unmarshal([]byte(cachedTweets), &tweets)
		if err == nil {
			log.Println("Cache hit")
			return tweets, nil
		}
	}
	//log.Println("simulating long query")
	//time.Sleep(5 * time.Second)
	query := `SELECT id, title, content, topic, created_at FROM tweets`

	rows, err := pg.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*domain.Tweet
	for rows.Next() {
		var tweet domain.Tweet
		err = rows.Scan(
			&tweet.ID,
			&tweet.Title,
			&tweet.Content,
			&tweet.Topic,
			&tweet.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Save the data to the cache
	tweetsJson, err := json.Marshal(tweets)
	if err == nil {
		pg.RedisClient.Set(ctx, cacheKey, tweetsJson, pg.CacheTTL)
		log.Println("Cache miss")
	}

	return tweets, nil
}

func (pg *repository) Update(in *domain.Tweet) (*domain.Tweet, error) {
	query := `
			UPDATE tweets 
			SET title = $1, content = $2, topic = $3, updated_at = NOW()
			WHERE id = $4
			RETURNING id, title, content, topic, updated_at`

	args := []interface{}{in.Title, in.Content, in.Topic, in.ID}
	err := pg.Db.QueryRow(context.Background(), query, args...).Scan(
		&in.ID,
		&in.Title,
		&in.Content,
		&in.Topic,
		&in.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, domain.ErrRecordNotFoundX
		default:
			return nil, err
		}
	}

	// Delete cache
	if err := pg.RedisClient.Del(context.Background(), "tweets:list").Err(); err != nil {
		return nil, fmt.Errorf("failed to delete cache: %w", err)

	}

	// Rebuild the cache
	if err := pg.RebuildCache(); err != nil {
		return nil, fmt.Errorf("failed to rebuild cache: %w", err)
	}

	return in, err
}

func (pg *repository) Delete(id int) error {
	query := `DELETE FROM tweets WHERE id = $1`

	result, err := pg.Db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrRecordNotFoundX
	}

	// Delete cache
	if err := pg.RedisClient.Del(context.Background(), "tweets:list").Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	// Rebuild the cache
	if err := pg.RebuildCache(); err != nil {
		return fmt.Errorf("failed to rebuild cache: %w", err)
	}

	return nil
}

func (pg *repository) GetUserTweets(id int) ([]*domain.Tweet, error) {
	query := `
			SELECT * FROM tweets
			WHERE user_id = $1`

	rows, err := pg.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*domain.Tweet
	for rows.Next() {
		var tweet *domain.Tweet
		err = rows.Scan(
			&tweet.ID,
			&tweet.Title,
			&tweet.Content,
			&tweet.Topic,
			&tweet.UserId,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}
