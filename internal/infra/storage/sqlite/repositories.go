package sqlite

import (
	"database/sql"

	"github.com/arnald/forum/internal/domain/categories"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/storage/sqlite/category"
	users "github.com/arnald/forum/internal/infra/storage/sqlite/user"
	"github.com/arnald/forum/internal/infra/storage/topics"
)

type Repositories struct {
	UserRepo     user.Repository
	CategoryRepo categories.Repository
	TopicRepo    topic.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:     users.NewRepo(db),
		CategoryRepo: category.NewRepo(db),
		TopicRepo:    topics.NewRepo(db),
	}
}
