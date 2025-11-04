package getalltopics

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	topicQueries "github.com/arnald/forum/internal/app/topics/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/topic"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type ResponseModel struct {
	Filters    map[string]interface{}   `json:"filters"`
	Topics     []topic.Topic            `json:"topics"`
	Pagination helpers.PaginationParams `json:"pagination"`
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

	var userID *string
	user := middleware.GetUserFromContext(r)
	if user != nil {
		userID = &user.ID
	}

	params := helpers.NewURLParams(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	pagination := params.GetPagination()

	orderBy := params.GetQueryStringOr("order_by", "created_at")
	order := params.GetQueryStringOr("order", "desc")
	filter := params.GetQueryStringOr("search", "")
	categoryID := params.GetQueryIntOr("category", 0)

	val := validator.New()

	validator.ValidateGetAllTopics(val, &struct {
		OrderBy    string
		Order      string
		Search     string
		CategoryID int
	}{
		OrderBy:    orderBy,
		Order:      order,
		Search:     filter,
		CategoryID: categoryID,
	})

	if !val.Valid() {
		h.Logger.PrintError(logger.ErrValidationFailed, val.Errors)
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			val.ToStringErrors(),
		)
		return
	}
	topics, totalCount, err := h.UserServices.UserServices.Queries.GetAllTopics.Handle(ctx, topicQueries.GetAllTopicsRequest{
		Page:       pagination.Page,
		Size:       pagination.Limit,
		Offset:     pagination.Offset,
		OrderBy:    orderBy,
		Order:      order,
		Filter:     filter,
		CategoryID: categoryID,
		UserID:     userID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get topics")
		return
	}

	totalPages := (totalCount + pagination.Limit - 1) / pagination.Limit

	paginationMeta := map[string]interface{}{
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total":       totalCount,
		"total_pages": totalPages,
		"has_next":    pagination.Page < totalPages,
		"has_prev":    pagination.Page > 1,
		"next_page":   nil,
		"prev_page":   nil,
	}

	if pagination.Page < totalPages {
		paginationMeta["next_page"] = pagination.Page + 1
	}
	if pagination.Page > 1 {
		paginationMeta["prev_page"] = pagination.Page - 1
	}

	appliedFilters := map[string]interface{}{
		"search":   filter,
		"order_by": orderBy,
		"order":    order,
	}

	response := map[string]interface{}{
		"topics":     topics,
		"pagination": paginationMeta,
		"filters":    appliedFilters,
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)

	h.Logger.PrintInfo("Topics retrieved successfully", map[string]string{
		"page":  strconv.Itoa(pagination.Page),
		"count": strconv.Itoa(len(topics)),
		"total": strconv.Itoa(totalCount),
	})
}
