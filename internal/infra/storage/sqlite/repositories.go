package sqlite

import (
	"database/sql"

	"github.com/arnald/forum/internal/domain/activity"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/notification"
	"github.com/arnald/forum/internal/domain/oauth"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/domain/vote"
	activities "github.com/arnald/forum/internal/infra/storage/sqlite/activity"
	"github.com/arnald/forum/internal/infra/storage/sqlite/categories"
	"github.com/arnald/forum/internal/infra/storage/sqlite/comments"
	oauthrepo "github.com/arnald/forum/internal/infra/storage/sqlite/oauth"
	"github.com/arnald/forum/internal/infra/storage/sqlite/topics"
	"github.com/arnald/forum/internal/infra/storage/sqlite/users"
	"github.com/arnald/forum/internal/infra/storage/sqlite/votes"
)

type Repositories struct {
	UserRepo         user.Repository
	CategoryRepo     category.Repository
	TopicRepo        topic.Repository
	CommentRepo      comment.Repository
	VoteRepo         vote.Repository
	NotificationRepo notification.Repository
	OauthRepo        oauth.Repository
	ActivityRepo     activity.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo:     users.NewRepo(db),
		CategoryRepo: categories.NewRepo(db),
		TopicRepo:    topics.NewRepo(db),
		CommentRepo:  comments.NewRepo(db),
		VoteRepo:     votes.NewRepo(db),
		OauthRepo:    oauthrepo.NewOAuthRepository(db),
		ActivityRepo: activities.NewRepo(db),
	}
}
