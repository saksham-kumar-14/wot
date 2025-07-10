package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	About     string    `json:"about"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) GetByID(ctx context.Context, userID int) (*User, error) {
	query := `SELECT username, email, about
	FROM users
	WHERE id = $1`

	// query timeout
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user User
	user.ID = int64(userID)
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Username,
		&user.Email,
		&user.About,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	// query timeout
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	return err
}
