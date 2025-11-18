package oauthlogin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	oauthservice "github.com/arnald/forum/internal/app/oauth"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	oathstate "github.com/arnald/forum/internal/pkg/oAuth"
)

type GitHubHandler struct {
	config         *config.ServerConfig
	loginService   *oauthservice.GitHubLoginService
	stateManager   *oathstate.StateManager
	sessionManager session.Manager
	logger         logger.Logger
}

func NewGitHubHandler(
	config *config.ServerConfig,
	loginService *oauthservice.GitHubLoginService,
	stateManager *oathstate.StateManager,
	sessionManager session.Manager,
	logger logger.Logger,
) *GitHubHandler {
	return &GitHubHandler{
		config:         config,
		loginService:   loginService,
		stateManager:   stateManager,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (h *GitHubHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := h.stateManager.Generate()
	if err != nil {
		h.logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	params := url.Values{}
	params.Add("client_id", h.config.OAuth.GitHub.ClientID)
	params.Add("redirect_uri", h.config.OAuth.GitHub.RedirectURL)
	params.Add("scope", strings.Join(h.config.OAuth.GitHub.Scopes, " "))
	params.Add("state", state)

	authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params.Encode())

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

}
