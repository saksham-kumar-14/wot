package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Likes     int64     `json:"likes"`
	Dislikes  int64     `json:"dislikes"`
	UserId    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
			VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserId,
		pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	return err
}

func (s *PostsStore) GetByID(ctx context.Context, postID int) (*Post, error) {
	query := `SELECT id, user_id, title, content, tags, likes, dislikes, created_at, updated_at
	FROM posts
	WHERE id = $1`

	var post Post
	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.Likes,
		&post.Dislikes,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostsStore) DeleteByID(ctx context.Context, postID int) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) PatchByID(ctx context.Context, postData *Post) error {
	query := `UPDATE posts
			  SET title = $1, content = $2, tags = $3
			  WHERE id = $4`

	result, err := s.db.ExecContext(ctx, query, postData.Title, postData.Content, pq.Array(postData.Tags), postData.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
