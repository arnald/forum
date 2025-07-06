package health

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arnald/forum/internal/app/health/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request method %v\n", r.Method)
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
