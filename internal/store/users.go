package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	About     string    `json:"about"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) GetByID(ctx context.Context, userID int) (*User, error) {
	query := `SELECT username, email, about
	FROM users
	WHERE id = $1`

	// query timeout
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user User
	user.ID = int64(userID)
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Username,
		&user.Email,
		&user.About,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	// query timeout
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password.hash,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key`:
			return DuplicateEmailErr
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key`:
			return DuplicateUsernameErr
		default:
			return err
		}
	}

	return nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, inviteExp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(inviteExp))
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) Activate(ctx context.Context, token string) error {
	// find user with token -> activate user -> clean up invitation
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// find user
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// activate user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// clean-up invitation
		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UsersStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `SELECT usr.id, usr.username, usr.email, usr.created_at, usr.is_active
	FROM users usr
	JOIN user_invitations ui ON usr.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > $2`

	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
