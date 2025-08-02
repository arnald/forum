package userlogin

import (
	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
)

type Handler struct {
	UserServices   app.Services
	SessionManager user.SessionManager
	Config         *config.ServerConfig
}

func NewHandler(config *config.ServerConfig, app app.Services, sm user.SessionManager) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
	}
}
