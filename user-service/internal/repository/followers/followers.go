package followers

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	Db *pgxpool.Pool
}

func NewFollowersRepo(db *pgxpool.Pool) *repo {
	return &repo{
		Db: db,
	}
}

func (pg *repo) Follow(followerID, followedID int) error {
	query := `INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2)`
	_, err := pg.Db.Exec(context.Background(), query, followerID, followedID)
	return err
}

func (pg *repo) Unfollow(followerID, followedID int) error {
	query := `DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2`
	result, err := pg.Db.Exec(context.Background(), query, followerID, followedID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrRecordNotFound
	}
	return nil
}

func (pg *repo) IsFollowing(followerID, followedID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND followed_id = $2)`
	var exists bool
	err := pg.Db.QueryRow(context.Background(), query, followerID, followedID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (pg *repo) GetFollowers(userID int) ([]*domain.User, error) {
	var users []*domain.User
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.username
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.followed_id = $1`

	rows, err := pg.Db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (pg *repo) GetFollowing(userID int) ([]*domain.User, error) {
	var users []*domain.User
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.username
		FROM users u
		JOIN followers f ON u.id = f.followed_id
		WHERE f.follower_id = $1`

	rows, err := pg.Db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
