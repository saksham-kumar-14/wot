package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostId    int64  `json:"post_id"`
	UserId    int64  `json:"user_id"`
	Content   string `json:"content"`
	Likes     int64  `json:"likes"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) GetCommentsHandler(ctx context.Context, postId int) ([]Comment, error) {
	query := `SELECT c.id, c.post_id, c.content, c.created_at, users.username, users.id
			  FROM comments c
			  JOIN users ON users.id = c.user_id
			  WHERE c.post_id = $1
			  ORDER BY c.created_at DESC;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}

		err := rows.Scan(
			&c.ID,
			&c.PostId,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (s *CommentsStore) CreateComment(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content)
			VALUES ($1, $2, $3) RETURNING id, created_at`

	// query timeout
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.PostId, comment.UserId,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	return err
}
