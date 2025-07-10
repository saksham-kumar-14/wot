package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Friend struct {
	UserID    int64  `json:"user_id"`
	FriendID  int64  `json:"friend_id"`
	CreatedAt string `json:"created_at"`
}

type FriendStore struct {
	db *sql.DB
}

func (s *FriendStore) Friend(ctx context.Context, userID int, friendID int) error {
	query := `
		INSERT INTO friends (user_id, friend_id) VALUES($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, friendID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrAlreadyExists
		}
	}

	return err
}

func (s *FriendStore) Unfriend(ctx context.Context, userID int, friendID int) error {
	query := `
		DELETE FROM friends WHERE user_id=$1 AND friend_id=$2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, friendID)
	return err
}
