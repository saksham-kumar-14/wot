package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("No document found")
)

type postDataType struct {
	Title   string
	Content string
}

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int) (*Post, error)
		DeleteByID(context.Context, int) error
		PatchByID(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		GetCommentsHandler(context.Context, int) ([]Comment, error)
	}
}

func NewDbStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
