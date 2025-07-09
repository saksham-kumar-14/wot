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
		JOIN users on users.id = c.user_id
		WHERE c.post_id=$1
		ORDER BY c.created_at DESC;`

	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostId, &c.UserId, &c.Content, &c.CreatedAt, &c.User.Username,
			&c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil

}
