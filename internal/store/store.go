package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("No document found")
	ErrAlreadyExists = errors.New("Resource already exsits")
	QueryTimeout     = time.Second * 5
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
		GetByID(context.Context, int) (*User, error)
	}
	Comments interface {
		GetCommentsHandler(context.Context, int) ([]Comment, error)
		CreateComment(context.Context, *Comment) error
	}
	Friends interface {
		Friend(context.Context, int, int) error
		Unfriend(context.Context, int, int) error
	}
}

func NewDbStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
		Friends:  &FriendStore{db},
	}
}
