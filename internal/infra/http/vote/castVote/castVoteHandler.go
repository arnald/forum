package castvote

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/arnald/forum/internal/app"
	votecommands "github.com/arnald/forum/internal/app/votes/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	TopicID      *int `json:"topicId"`
	CommentID    *int `json:"commentId,omitempty"`
	ReactionType int  `json:"reactionType"`
}

type ResponseModel struct {
	Message    string           `json:"message"`
	VoteCounts *vote.VoteCounts `json:"voteCounts"`
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

func (h *Handler) CastVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(
			w,
			http.StatusMethodNotAllowed,
			"Invalid request method",
		)
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		h.Logger.PrintError(logger.ErrUserNotFoundInContext, nil)
		helpers.RespondWithError(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var req RequestModel
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid JSON",
		)
	}

	targer := vote.VoteTarget{
		TopicID:   req.TopicID,
		CommentID: req.CommentID,
	}

	err = h.Services.UserServices.Commands.CastVote.Handle(ctx, votecommands.CastVoteRequest{
		UserID:       user.ID,
		Target:       targer,
		ReactionType: req.ReactionType,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"failed to cast vote",
		)
		return
	}

	ResponseModel := ResponseModel{
		Message: "Vote cast successfully",
	}

	helpers.RespondWithJSON(w,
		http.StatusOK,
		nil,
		ResponseModel,
	)
}
