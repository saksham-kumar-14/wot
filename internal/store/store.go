package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("No document found")
	ErrAlreadyExists     = errors.New("Resource already exsits")
	DuplicateEmailErr    = errors.New("Email already exists")
	DuplicateUsernameErr = errors.New("Username already exsits")
	QueryTimeout         = time.Second * 5
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
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
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

func withTx(db *sql.DB, ctx context.Context, function func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := function(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
