package repository

import (
	"MussaShaukenov/twitter-clone-go/internal/tweet/domain"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	Db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *postgres {
	return &postgres{
		Db: db,
	}
}

func (pg *postgres) Insert(in *domain.Tweet) error {
	query := `
				INSERT INTO tweets (title, content, topic) 
				VALUES ($1, $2, $3)
				RETURNING id, created_at`

	args := []interface{}{in.Title, in.Content, in.Topic}
	return pg.Db.QueryRow(context.Background(), query, args...).
		Scan(&in.ID, &in.CreatedAt)
}

func (pg *postgres) Get(id int64) (*domain.Tweet, error) {
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

func (pg *postgres) List() ([]*domain.Tweet, error) {
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

	return tweets, nil
}

func (pg *postgres) Update(in *domain.Tweet) (*domain.Tweet, error) {
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

	return in, err
}

func (pg *postgres) Delete(id int) error {
	query := `DELETE FROM tweets WHERE id = $1`

	result, err := pg.Db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrRecordNotFoundX
	}
	return nil
}
