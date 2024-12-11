package users

import (
	"MussaShaukenov/twitter-clone-go/user-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type repository struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewUsersRepo(db *pgxpool.Pool, logger *zap.SugaredLogger) *repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (repo *repository) Insert(in *domain.User) error {
	repo.logger.Infow("Inserting user into repository", "user", in)
	query := `
			INSERT INTO users (first_name, last_name, email, username, password, age) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at`

	args := []interface{}{in.FirstName, in.LastName, in.Email, in.Username, in.Password, in.Age}
	err := repo.db.QueryRow(context.Background(), query, args...).Scan(&in.ID, &in.CreatedAt)
	if err != nil {
		repo.logger.Errorw("Failed to insert user", "error", err)
		return err
	}
	return nil
}

func (repo *repository) GetByID(id int) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username 
		FROM users 
		WHERE id = $1`

	var user domain.User
	err := repo.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.logger.Errorw("Failed to get user by id", "error", err)
			return nil, domain.ErrRecordNotFound
		}
		repo.logger.Errorw("Failed to get user by id", "error", err)
		return nil, err
	}
	return &user, nil
}

func (repo *repository) GetUserEmail(id int) (string, error) {
	query := `
		SELECT email FROM users
		WHERE id = $1`

	var email string
	err := repo.db.QueryRow(context.Background(), query, id).Scan(&email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			repo.logger.Errorw("Failed to get user email by id", "error", err)
			return "", domain.ErrRecordNotFound
		default:
			repo.logger.Errorw("Failed to get user email by id", "error", err)
			return "", err
		}
	}

	return email, nil
}

func (repo *repository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := repo.db.Exec(context.Background(), query, id)
	if err != nil {
		repo.logger.Errorw("Failed to delete user", "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		repo.logger.Errorw("Failed to delete user", "error", domain.ErrRecordNotFound)
		return domain.ErrRecordNotFound
	}
	return nil
}

func (repo *repository) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password 
		FROM users 
		WHERE username = $1`

	var user domain.User
	err := repo.db.QueryRow(context.Background(), query, username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.logger.Errorw("Failed to get user by username", "error", err)
			return nil, domain.ErrRecordNotFound
		}
		repo.logger.Errorw("Failed to get user by username", "error", err)
		return nil, err
	}
	return &user, nil
}

func (repo *repository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, first_name, last_name, email, username, password FROM users WHERE email = $1`
	err := repo.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			repo.logger.Errorw("Failed to get user by email", "error", err)
			return nil, domain.ErrRecordNotFound
		default:
			repo.logger.Errorw("Failed to get user by email", "error", err)
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}
	return &user, nil
}

func (repo *repository) IsFirstLogin(userId int) (bool, error) {
	var user domain.User
	query := `SELECT id, is_first_login FROM users WHERE id = $1`
	err := repo.db.QueryRow(context.Background(), query, userId).Scan(&user.ID, &user.IsFirstLogin)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			repo.logger.Errorw("Failed to get first login status", "error", err)
			return false, domain.ErrRecordNotFound
		default:
			repo.logger.Errorw("Failed to get first login status", "error", err)
			return false, fmt.Errorf("failed to get first login status: %w", err)
		}
	}

	return !user.IsFirstLogin, nil
}

func (repo *repository) List() ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT id, first_name, last_name, email, username FROM users`
	rows, err := repo.db.Query(context.Background(), query)
	if err != nil {
		repo.logger.Errorw("Failed to list users", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username)
		if err != nil {
			repo.logger.Errorw("Failed to scan user", "error", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
