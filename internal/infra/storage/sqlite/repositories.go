package sqlite

import (
	"database/sql"

	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/storage/sqlite/categories"
	"github.com/arnald/forum/internal/infra/storage/sqlite/topics"
	"github.com/arnald/forum/internal/infra/storage/sqlite/users"
)

type Repositories struct {
	UserRepo     user.Repository
	CategoryRepo category.Repository
	TopicRepo    topic.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:     users.NewRepo(db),
		CategoryRepo: categories.NewRepo(db),
		TopicRepo:    topics.NewRepo(db),
	}
}
