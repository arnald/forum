package infra

import (
	"database/sql"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/comment"
	"github.com/arnald/forum/internal/domain/post"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/infra/http"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
)

type Services struct {
	UserRepository     user.Repository
	PostRepository     post.Repository
	CommentRepository  comment.Repository
	VoteRepository     vote.Repository
	CategoryRepository category.Repository
	Server             *http.Server
}

func NewInfraProviders(db *sql.DB) Services {
	repo := sqlite.NewRepo(db)
	return Services{
		UserRepository:     repo,
		PostRepository:     repo,
		CommentRepository:  repo,
		VoteRepository:     repo,
		CategoryRepository: repo,
	}
}

func NewHTTPServer(cfg *config.ServerConfig, db *sql.DB, logger logger.Logger, appServices app.Services) *http.Server {
	return http.NewServer(cfg, db, logger, appServices)
}
