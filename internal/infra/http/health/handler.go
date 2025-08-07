package health

import (
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app/health/queries"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	Logger logger.Logger
}

func NewHandler(logger logger.Logger) *Handler {
	return &Handler{
		Logger: logger,
	}
}

func (h Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")

		return
	}

	response := queries.HealthResponse{
		Status:    queries.StatusUp,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		response,
	)
}
