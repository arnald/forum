package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/user"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

// TODO: retrieves all users from the repository.
func (r Repo) GetAll(_ context.Context) ([]user.User, error) {
	return nil, nil
}

func (r Repo) UserRegister(ctx context.Context, user *user.User) error {
	query := `
	INSERT INTO users (username, password_hash, email, id)
	VALUES (?, ?, ?, ?)`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	_, err = r.DB.ExecContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
		user.ID,
	)

	mapErr := MapSQLiteError(err)
	if mapErr != nil {
		return mapErr
	}

	return nil
}

func (r Repo) GetUserByIdentifier(ctx context.Context, identifier string) (*user.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, avatar_url
	FROM users
	WHERE email = ? OR username = ?
	`
	var user user.User
	err := r.DB.QueryRowContext(ctx, query, identifier, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.AvatarURL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user with identifier %s not found: %w", identifier, ErrUserNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by identifier: %w", err)
	}

	return &user, nil
}

func (r Repo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
	SELECT id, username, password_hash
	FROM users
	WHERE email = ?
	`
	var user user.User
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user with email %s not found: %w", email, ErrUserNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r Repo) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, avatar_url
	FROM users
	WHERE username = ?
	`
	var user user.User
	err := r.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.AvatarURL,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user with username %s not found: %w", username, ErrUserNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}
