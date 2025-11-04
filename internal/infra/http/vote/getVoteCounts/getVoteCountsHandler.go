package getvotecounts

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	votequeries "github.com/arnald/forum/internal/app/votes/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/vote"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Response struct {
	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
	Score     int `json:"score"`
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

func (h *Handler) GetCounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var topicID *int
	topicValue, err := helpers.GetQueryInt(r, "topic_id")
	if err == nil {
		topicID = &topicValue
	}
	var commentID *int
	commentValue, err := helpers.GetQueryInt(r, "comment_id")
	if err == nil {
		commentID = &commentValue
	}

	if topicID != nil && commentID != nil {
		h.Logger.PrintError(logger.ErrBothIDsProvided, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "Provide either topic_id OR comment_id, not both")
		return
	}

	if topicID == nil && commentID == nil {
		h.Logger.PrintError(logger.ErrNeitherIDProvided, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "Either topic_id or comment_id is required")
		return
	}

	Target := vote.Target{
		TopicID:   topicID,
		CommentID: commentID,
	}

	Counts, err := h.Services.UserServices.Queries.GetCounts.Handle(ctx, votequeries.GetCountsRequest{
		Target: Target,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		http.Error(w, "Failed to get vote counts", http.StatusInternalServerError)
		return
	}

	response := Response{
		Upvotes:   Counts.Upvotes,
		Downvotes: Counts.DownVotes,
		Score:     Counts.Score,
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		response,
	)
	logData := map[string]string{}

	if topicID != nil {
		logData["topicID"] = strconv.Itoa(*topicID)
	} else {
		logData["topicID"] = "nil"
	}

	if commentID != nil {
		logData["commentID"] = strconv.Itoa(*commentID)
	} else {
		logData["commentID"] = "nil"
	}

	h.Logger.PrintInfo("Vote counts retrieved successfully", logData)
}
