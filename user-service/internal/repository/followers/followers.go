package followers

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type repository struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewFollowersRepo(db *pgxpool.Pool, logger *zap.SugaredLogger) *repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (repo *repository) Follow(followerID, followedID int) error {
	query := `INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2)`
	_, err := repo.db.Exec(context.Background(), query, followerID, followedID)
	if err != nil {
		repo.logger.Errorw("Failed to follow", "followerID", followerID, "followedID", followedID, "error", err)
		return err
	}
	return nil
}

func (repo *repository) Unfollow(followerID, followedID int) error {
	query := `DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2`
	result, err := repo.db.Exec(context.Background(), query, followerID, followedID)
	if err != nil {
		repo.logger.Errorw("Failed to unfollow", "followerID", followerID, "followedID", followedID, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		repo.logger.Warnw("Failed to unfollow", "followerID", followerID, "followedID", followedID, "error", domain.ErrRecordNotFound)
		return domain.ErrRecordNotFound
	}
	return nil
}

func (repo *repository) IsFollowing(followerID, followedID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND followed_id = $2)`
	var exists bool
	err := repo.db.QueryRow(context.Background(), query, followerID, followedID).Scan(&exists)
	if err != nil {
		repo.logger.Errorw("Failed to check if following", "followerID", followerID, "followedID", followedID, "error", err)
		return false, err
	}
	return exists, nil
}

func (repo *repository) GetFollowers(userID int) ([]*domain.User, error) {
	var users []*domain.User
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.username
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.followed_id = $1`

	rows, err := repo.db.Query(context.Background(), query, userID)
	if err != nil {
		repo.logger.Errorw("Failed to get followers", "userID", userID, "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			repo.logger.Errorw("Failed to scan row", "error", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (repo *repository) GetFollowing(userID int) ([]*domain.User, error) {
	var users []*domain.User
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.username
		FROM users u
		JOIN followers f ON u.id = f.followed_id
		WHERE f.follower_id = $1`

	rows, err := repo.db.Query(context.Background(), query, userID)
	if err != nil {
		repo.logger.Errorw("Failed to get following", "userID", userID, "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			repo.logger.Errorw("Failed to scan row", "error", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
