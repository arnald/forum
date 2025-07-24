package sqlite

import (
	"context"
	"database/sql"
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

func (r Repo) CreateSession(session *user.Session) error {
	ctx := context.TODO()

	query := `
	INSERT INTO sessions (token, user_id, expires_at)
	VALUES (?, ?, ?)`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		session.Token,
		session.UserID,
		session.Expiry.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
	SELECT id, username, password_hash, email
	FROM users
	WHERE email = ?`

	row := r.DB.QueryRowContext(ctx, query, email)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return &u, nil
}
