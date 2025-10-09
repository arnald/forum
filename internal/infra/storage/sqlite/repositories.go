package sqlite

import (
	"database/sql"

	"github.com/arnald/forum/internal/domain/categories"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/storage/sqlite/category"
)

type Repositories struct {
	UserRepo     user.Repository
	CategoryRepo categories.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:     NewRepo(db),
		CategoryRepo: category.NewRepo(db),
	}
}
