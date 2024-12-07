package tags

import (
	"MussaShaukenov/twitter-clone-go/tweet-service/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	Db *pgxpool.Pool
}

func NewTagsRepository(db *pgxpool.Pool) *repository {
	return &repository{
		Db: db,
	}
}

func (pg *repository) AddTag(tweetId int64, tagId int64) error {
	query := `
			INSERT INTO tweet_tags (tweet_id, tag_id)
			VALUES ($1, $2)`

	_, err := pg.Db.Exec(context.Background(), query, tweetId, tagId)
	if err != nil {
		fmt.Errorf("error adding tag to tweet: %v", err)
		return err
	}
	return nil
}

func (pg *repository) GetTweetTags(tweetId int64) ([]*domain.Tag, error) {
	query := `
			SELECT t.id, t.name FROM tags t
			JOIN tweet_tags tt ON t.id = tt.tag_id
			WHERE tt.tweet_id = $1`

	rows, err := pg.Db.Query(context.Background(), query, tweetId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (pg *repository) ListTags() ([]*domain.Tag, error) {
	query := `SELECT id, name FROM tags`

	rows, err := pg.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
