package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Likes     int64    `json:"likes"`
	Dislikes  int64    `json:"dislikes"`
	UserID    int64    `json:"user_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Tags      []string `json:"tags"`
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
			VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, post.Content, post, post.Title, post.UserID,
		pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	return err
}
