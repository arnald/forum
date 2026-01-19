package getuseractivity

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	activityQueries "github.com/arnald/forum/internal/app/activities/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/activity"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type ResponseModel struct {
	CreatedTopics    []activity.TopicActivity       `json:"createdTopics"`
	LikedTopics      []activity.TopicActivity       `json:"likedTopics"`
	DislikedTopics   []activity.TopicActivity       `json:"dislikedTopics"`
	LikedComments    []activity.CommentVoteActivity `json:"likedComments"`
	DislikedComments []activity.CommentVoteActivity `json:"dislikedComments"`
	UserComments     []activity.CommentActivity     `json:"userComments"`
}
type Handler struct {
	Services app.Services
	Config   *config.ServerConfig
	Logger   logger.Logger
}

func NewHandler(services app.Services, config *config.ServerConfig, logger logger.Logger) *Handler {
	return &Handler{
		Services: services,
		Config:   config,
		Logger:   logger,
	}
}

func (h *Handler) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	user := middleware.GetUserFromContext(r)
	if user == nil {
		h.Logger.PrintError(logger.ErrUserNotFoundInContext, nil)
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	activity, err := h.Services.UserServices.Queries.GetUserActivity.Handle(ctx, activityQueries.GetUserActivityRequest{
		UserID: user.ID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get user activity")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, activity)
	h.Logger.PrintInfo("User activity retrieved successfully", map[string]string{"userID": user.ID})
}
