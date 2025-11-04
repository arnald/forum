package markallasread

import (
	"net/http"

	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/notifications"
)

type Handler struct {
	service *notifications.NotificationService
}

func NewHandler(service *notifications.NotificationService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	if user.ID == "" {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	userID := user.ID

	err := h.service.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		http.Error(
			w,
			"failed to mark all as read",
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}
