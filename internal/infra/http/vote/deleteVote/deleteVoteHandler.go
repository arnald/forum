package deletevote

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/arnald/forum/internal/app"
	votecommands "github.com/arnald/forum/internal/app/votes/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type request struct {
	TopicID   *int   `json:"topicId,omitempty"`
	CommentID *int   `json:"commentId,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

type Response struct {
	Message string
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

func (h *Handler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(
			w,
			http.StatusMethodNotAllowed,
			"invalid method",
		)
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		h.Logger.PrintError(logger.ErrUserNotFoundInContext, nil)
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to read requst body")
	}

	voteToDelete := &request{}

	err = json.Unmarshal(data, voteToDelete)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to unmarshal request")
	}

	err = h.Services.UserServices.Commands.DeleteVote.Handle(ctx, votecommands.DeleteVoteRequest{
		UserID:    user.ID,
		TopicID:   voteToDelete.TopicID,
		CommentID: voteToDelete.CommentID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"failed to delete vote",
		)
		return
	}

	voteResponse := Response{
		Message: "Vote deleted successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		voteResponse,
	)

	h.Logger.PrintInfo(
		"Vote deleted successfully",
		nil,
	)
}
