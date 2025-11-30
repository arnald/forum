package logout

import (
	"net/http"

	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	sessionManager session.Manager
	logger         logger.Logger
}

func NewHandler(sessionManager session.Manager, logger logger.Logger) *Handler {
	return &Handler{
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// Logout deletes the user's session from the database.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Get user from context (set by Required auth middleware)
	user := middleware.GetUserFromContext(r)
	if user == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get session token from cookie
	sessionToken, _ := middleware.GetTokensFromRequest(r)
	if sessionToken == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "No session found")
		return
	}

	// Delete the session from database
	err := h.sessionManager.DeleteSession(sessionToken)
	if err != nil {
		h.logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	h.logger.PrintInfo("User logged out successfully", map[string]string{
		"userId": user.ID,
		"name":   user.Username,
	})

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]string{
		"message": "Logged out successfully",
	})
}
