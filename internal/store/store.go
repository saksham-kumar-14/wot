package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("no document found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")

	QueryTimeout = time.Second * 5
)

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
		GetByEmail(context.Context, string) (*User, error)
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
