package markasread

import (
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/infra/storage/notifications"
)

type Handler struct {
	service *notifications.NotificationService
}

func NewHandler(service *notifications.NotificationService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
	}

	if user.ID == "" {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
	}

	userID := user.ID

	notificationID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(
			w,
			"invalida notification ID",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.MarkAsRead(r.Context(), int(notificationID), userID)
	if err != nil {
		http.Error(
			w,
			"failed to mark as read",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
