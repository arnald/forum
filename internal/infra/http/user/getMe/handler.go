package getme

import (
	"net/http"

	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	logger logger.Logger
}

func NewHandler(logger logger.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

type GetMeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GetMe handler retrieves the current user from the session in the context.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Get user from context (set by middleware) using the helper function
	user := middleware.GetUserFromContext(r)
	if user == nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: User not found")
		return
	}

	response := GetMeResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)

	h.logger.PrintInfo(
		"User retrieved from session",
		map[string]string{
			"userId": user.ID,
			"name":   user.Username,
		},
	)
}
