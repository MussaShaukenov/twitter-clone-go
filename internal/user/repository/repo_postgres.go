package repository

import (
	"MussaShaukenov/twitter-clone-go/internal/user/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type postgresRepo struct {
	Db *pgxpool.Pool
}

func NewPostgres(db *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{
		Db: db,
	}
}

func (pg *postgresRepo) Insert(in *domain.User) error {
	log.Println("User: ", in.FirstName, in.LastName, in.Email, in.Username, in.Password)
	query := `
			INSERT INTO users (first_name, last_name, email, username, password) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

	args := []interface{}{in.FirstName, in.LastName, in.Email, in.Username, in.Password}
	return pg.Db.QueryRow(context.Background(), query, args...).Scan(&in.ID, &in.CreatedAt)
}

func (pg *postgresRepo) Get(id int) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username 
		FROM users 
		WHERE id = $1`

	var user domain.User

	err := pg.Db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (pg *postgresRepo) GetUserEmail(id int) (string, error) {
	query := `
		SELECT email FROM users
		WHERE id = $1`

	var email string
	err := pg.Db.QueryRow(context.Background(), query, id).Scan(&email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", domain.ErrRecordNotFound
		default:
			return "", err
		}
	}

	return email, nil
}

func (pg *postgresRepo) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := pg.Db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrRecordNotFound
	}
	return nil
}

func (pg *postgresRepo) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password 
		FROM users 
		WHERE username = $1`

	var user domain.User
	err := pg.Db.QueryRow(context.Background(), query, username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, first_name, last_name, email, username, password FROM users WHERE email = $1`
	err := r.Db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}
	return &user, nil
}

func (pg *postgresRepo) IsFirstLogin(userId int) (bool, error) {
	var user domain.User
	query := `SELECT id, is_first_login FROM users WHERE id = $1`
	err := pg.Db.QueryRow(context.Background(), query, userId).Scan(&user.ID, &user.IsFirstLogin)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, domain.ErrRecordNotFound
		default:
			return false, fmt.Errorf("failed to get first login status: %w", err)
		}
	}

	return !user.IsFirstLogin, nil

}
