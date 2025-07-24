package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/arnald/forum/internal/domain/user"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo() Repo {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	return Repo{
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

func (r Repo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, username, password_hash, email FROM users WHERE email = ?")
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var u user.User
	err = row.Scan(&u.ID, &u.Username, &u.Password, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return &u, nil
}
