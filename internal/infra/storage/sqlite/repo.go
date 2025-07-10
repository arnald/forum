package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/go-sql-driver/mysql"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo() Repo {
	db, err := sql.Open("mysql", "georgeoik:123@tcp(localhost:3306)/forum")
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
	query := `
	INSERT INTO users (username, password, email, ID, created_at)
	VALUES (?, ?, ?, ?, ?)`

	_, err := r.DB.ExecContext(context.TODO(), query, user.Username, encryptedPass, user.Email, user.ID.String(), user.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		mysqlErr := &mysql.MySQLError{}
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == uniqueConstraintViolationErrorCode {
				return ErrDuplicateEmail
			}
			return fmt.Errorf("mysql error %d: %s", mysqlErr.Number, mysqlErr.Message)
		}
		return err
	}

	return nil
}

func (r Repo) CreateSession(session *user.Session) error {
	query := `
        INSERT INTO sessions (token, user_id, expiry, ip_address)
        VALUES (?, ?, ?, ?)`

	_, err := r.DB.ExecContext(
		context.TODO(),
		query,
		session.Token,
		session.UserID,
		session.Expiry,
		session.IPAddress,
	)
	if err != nil {
		return err
	}

	return nil
}
