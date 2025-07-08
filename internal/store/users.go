package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	About     string    `json:"about"`
	CreatedAt time.Time `json:"created_at"`
	Friends   []string  `json:"friends"`
	FriendsOf []string  `json:"friends_of"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password, about, friends, friends_of)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.About,
		pq.Array(user.Friends),
		pq.Array(user.FriendsOf),
	).Scan(&user.ID, &user.CreatedAt)

	return err
}
