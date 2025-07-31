package infra

import (
	"database/sql"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/http"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
)

type Services struct {
	UserRepository user.Repository
	Server         *http.Server
}

func NewInfraProviders(DB *sql.DB) Services {
	return Services{
		UserRepository: sqlite.NewRepo(DB),
	}
}

func NewHTTPServer(cfg *config.ServerConfig, db *sql.DB, appServices app.Services) *http.Server {
	return http.NewServer(cfg, db, appServices)
}
