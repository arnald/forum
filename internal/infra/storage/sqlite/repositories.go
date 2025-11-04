package sqlite

import (
	"database/sql"

	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/infra/storage/sqlite/categories"
	"github.com/arnald/forum/internal/infra/storage/sqlite/comments"
	"github.com/arnald/forum/internal/infra/storage/sqlite/topics"
	"github.com/arnald/forum/internal/infra/storage/sqlite/users"
	"github.com/arnald/forum/internal/infra/storage/sqlite/votes"
)

type Repositories struct {
	UserRepo     user.Repository
	CategoryRepo category.Repository
	TopicRepo    topic.Repository
	CommentRepo  comment.Repository
	VoteRepo     vote.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:     users.NewRepo(db),
		CategoryRepo: categories.NewRepo(db),
		TopicRepo:    topics.NewRepo(db),
		CommentRepo:  comments.NewRepo(db),
		VoteRepo:     votes.NewRepo(db),
	}
}
