package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/mattn/go-sqlite3"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo() Repo {
	db, err := sql.Open("sqlite3", "./forum.db")
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

func (r Repo) UserRegister(user *user.User, encryptedPass []byte) error {
	ctx := context.TODO()

	query := `
	INSERT INTO users (username, password, email, id, created_at)
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmt.Close()

	_, err = r.DB.ExecContext(
		context.TODO(),
		query,
		user.Username,
		string(encryptedPass),
		user.Email,
		user.ID.String(),
		user.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				switch sqliteErr.ExtendedCode {
				case sqlite3.ErrConstraintUnique:
					return ErrDuplicateEmail
				case sqlite3.ErrConstraintPrimaryKey:
					return fmt.Errorf("user with this ID already exists")
				default:
					return fmt.Errorf("sqlite constraint error: %v", sqliteErr)
				}
			}
			return fmt.Errorf("sqlite error %d: %s", sqliteErr.Code, sqliteErr.Error())
		}
		return err
	}

	return nil
}

func (r Repo) CreateSession(session *user.Session) error {
	ctx := context.TODO()

	query := `
	INSERT INTO sessions (token, user_id, expiry, ip_address)
	VALUES (?, ?, ?, ?)`

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
		session.IPAddress,
	)
	if err != nil {
		return err
	}

	return nil
}
