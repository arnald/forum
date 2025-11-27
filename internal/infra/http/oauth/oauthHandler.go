package oauthlogin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	oauthservice "github.com/arnald/forum/internal/app/oauth"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	oauthpkg "github.com/arnald/forum/internal/pkg/oAuth"
)

type OAuthHandler struct {
	provider       oauthpkg.Provider
	config         *config.ServerConfig
	loginService   *oauthservice.OAuthService
	stateManager   *oauthpkg.StateManager
	sessionManager session.Manager
	logger         logger.Logger
}

func NewOAuthHandler(
	provider oauthpkg.Provider,
	config *config.ServerConfig,
	loginService *oauthservice.OAuthService,
	stateManager *oauthpkg.StateManager,
	sessionManager session.Manager,
	logger logger.Logger,
) *OAuthHandler {
	return &OAuthHandler{
		provider:       provider,
		config:         config,
		loginService:   loginService,
		stateManager:   stateManager,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	authURL := h.provider.GetAuthURL(state)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	errParam := r.URL.Query().Get("error")
	if errParam != "" {
		h.logger.PrintError(fmt.Errorf("%w: %s", ErrInParameters, errParam), nil)
		http.Error(
			w,
			"problem with oatuh, see logger",
			http.StatusInternalServerError,
		)
		return
	}

	if code == "" {
		h.logger.PrintError(ErrCodeMissing, nil)
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
		ctx,
		code,
		h.provider,
	)
	if err != nil {
		h.logger.PrintError(err, map[string]string{
			"action":   "oauth_login",
			"provider": h.provider.Name(),
		})
		http.Error(
			w,
			"error at github_login",
			http.StatusInternalServerError,
		)
		return
	}

	session, err := h.sessionManager.CreateSession(r.Context(), user.ID)
	if err != nil {
		h.logger.PrintError(err, nil)
		http.Error(
			w,
			"error at creating session",
			http.StatusInternalServerError,
		)
	}

	params := url.Values{}
	params.Add("access_token", session.AccessToken)
	params.Add("refresh_token", session.RefreshToken)

	frontendCallbackURL := fmt.Sprintf("%s?%s", h.config.OAuth.FrontendCallbackURL, params.Encode())

	http.Redirect(w, r, frontendCallbackURL, http.StatusTemporaryRedirect)

	h.logger.PrintInfo(
		"User logged in via "+h.provider.Name(),
		map[string]string{
			"user_id":  user.ID,
			"username": user.Username,
			"provider": h.provider.Name(),
		})
}
