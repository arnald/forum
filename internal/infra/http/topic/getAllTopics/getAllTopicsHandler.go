package getalltopics

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/arnald/forum/internal/app"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type RequestModel struct {
	OrderBy  string `json:"orderBy,omitempty"`
	Filter   string `json:"filter,omitempty"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type ResponseModel struct {
	Topics      []topic.Topic `json:"topics"`
	TotalCount  int           `json:"totalCount"`
	CurrentPage int           `json:"currentPage"`
	PageSize    int           `json:"pageSize"`
}

type Handler struct {
	UserServices app.Services
	Config       *config.ServerConfig
	Logger       logger.Logger
}

func NewHandler(userServices app.Services, config *config.ServerConfig, logger logger.Logger) *Handler {
	return &Handler{
		UserServices: userServices,
		Config:       config,
		Logger:       logger,
	}
}

func (h *Handler) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var req RequestModel
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v := validator.New()

	validator.ValidateGetAllTopics(v, &req)
	if !v.Valid() {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			v.ToStringErrors(),
		)

		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)
		return
	}

	topics, count, err := h.UserServices.UserServices.Queries.GetAllTopics.Handle(ctx, topicQueries.GetAllTopicsRequest{
		Page:    req.Page,
		Size:    req.PageSize,
		OrderBy: req.OrderBy,
		Filter:  req.Filter,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get topics")
		return
	}

	response := ResponseModel{
		Topics:      topics,
		TotalCount:  count,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
	}
	helpers.RespondWithJSON(w, http.StatusOK, nil, response)
}
