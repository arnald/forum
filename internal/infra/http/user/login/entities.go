package userlogin

import (
	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/logger"
)

type Handler struct {
	UserServices   app.Services
	SessionManager user.SessionManager
	Config         *config.ServerConfig
	Logger         logger.Logger
}

func NewHandler(config *config.ServerConfig, app app.Services, sm user.SessionManager, logger logger.Logger) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
		Logger:         logger,
	}
}
