package health

import (
	"net/http"
	"time"

	"github.com/arnald/forum/internal/app/health/queries"
	"github.com/arnald/forum/internal/domain/notification"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/storage/notifications"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	Logger        logger.Logger
	Notifications *notifications.NotificationService
}

func NewHandler(logger logger.Logger, notifications *notifications.NotificationService) *Handler {
	return &Handler{
		Logger:        logger,
		Notifications: notifications,
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
	userID := "df16d238-e4dd-4645-9101-54aed9c0fbf4"
	relatedID := "testUserID"
	testNotification := &notification.Notification{
		UserID:      userID,
		Type:        notification.NotificationTypeLike,
		Title:       "Test Notification",
		Message:     "This is a test notification from health check",
		RelatedType: "post",
		RelatedID:   relatedID,
		IsRead:      false,
	}
	err := h.Notifications.CreateNotification(r.Context(), testNotification)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create test notification")

		return
	}
	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		response,
	)
}
