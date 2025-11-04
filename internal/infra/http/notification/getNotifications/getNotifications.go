package getnotifications

import (
	"encoding/json"
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

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
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

	limit := 50
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	notifications, err := h.service.GetNotifications(
		r.Context(),
		userID,
		limit,
	)
	if err != nil {
		http.Error(
			w,
			"failed to fetch notifications",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(notifications)
	if err != nil {
		http.Error(
			w,
			"failed to encode notifications",
			http.StatusInternalServerError,
		)
		return
	}
}
