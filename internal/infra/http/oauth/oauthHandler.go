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

func (h *GitHubHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	errParam := r.URL.Query().Get("error")
	if errParam != "" {
		h.logger.PrintError(fmt.Errorf("github oauth error: %s", errParam), nil)
		http.Error(
			w,
			"problem with oatuh, see logger",
			http.StatusInternalServerError,
		)
		return
	}

	if code == "" {
		h.logger.PrintError(fmt.Errorf("no code in callback"), nil)
		http.Error(
			w,
			"no code in callback",
			http.StatusInternalServerError,
		)
		return
	}

	err := h.stateManager.Verify(state)
	if err != nil {
		h.logger.PrintError(err, nil)
		http.Error(
			w,
			"problem with oauth STATE, SEE LOGGER",
			http.StatusInternalServerError,
		)
		return
	}

	user, err := h.loginService.Login(
		r.Context(),
		code,
		h.config.OAuth.GitHub.ClientID,
		h.config.OAuth.GitHub.ClientSecret,
		h.config.OAuth.GitHub.RedirectURL,
	)

	if err != nil {
		h.logger.PrintError(err, map[string]string{"action": "github_login"})
		http.Error(
			w,
			"error at github_login",
			http.StatusInternalServerError,
		)
		return
	}

	sessionID, err := h.sessionManager.CreateSession(r.Context(), user.ID)
	if err != nil {
		h.logger.PrintError(err, nil)
		http.Error(
			w,
			"error at creating session",
			http.StatusInternalServerError,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: sessionID.AccessToken,
		Path:  "/",
	})

	h.logger.PrintInfo(
		"USER LOGGED IN VIA GITHUB",
		map[string]string{
			"user_id":  user.ID,
			"username": user.Username,
		})

}
