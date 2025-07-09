package sqlite

import (
	"context"
	"database/sql"
	"log"

	"github.com/arnald/forum/internal/domain/user"
	_ "github.com/go-sql-driver/mysql"
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

func (r Repo) UserRegister(user user.User) error {
	query := `
	INSERT INTO users (username, password, email, ID, created_at, role)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(query, user.Username, user.Password, user.Email, user.ID.String(), user.CreatedAt.Format("2006-01-02 15:04:05"), user.Role)
	if err != nil {
		return err
	}

	return nil

}
