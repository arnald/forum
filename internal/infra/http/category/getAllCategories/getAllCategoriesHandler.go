package getallcategories

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	categoryqueries "github.com/arnald/forum/internal/app/categories/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type ResponseModel struct {
	Filters    map[string]interface{}   `json:"filters"`
	Categories []category.Category      `json:"categories"`
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

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	params := helpers.NewURLParams(r)

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	pagination := params.GetPagination()

	orderBy := params.GetQueryStringOr("order_by", "created_at")
	order := params.GetQueryStringOr("order", "desc")
	filter := params.GetQueryStringOr("search", "")

	categories, totalCount, err := h.UserServices.UserServices.Queries.GettAllCategories.Handle(ctx, categoryqueries.GetAllCategoriesRequest{
		OrderBy: orderBy,
		Order:   order,
		Filter:  filter,
		Page:    pagination.Page,
		Size:    pagination.Limit,
		Offset:  pagination.Offset,
	})

	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to get categories",
		)
		return
	}

	totalPages := (totalCount + pagination.Limit - 1) / pagination.Limit

	paginationMeta := map[string]interface{}{
		"page":       pagination.Page,
		"limit":      pagination.Limit,
		"totalPages": totalPages,
		"totalItems": totalCount,
		"has_next":   pagination.Page < totalPages,
		"has_prev":   pagination.Page > 1,
		"next_page":  nil,
		"prev_page":  nil,
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
		"categories": categories,
		"pagination": paginationMeta,
		"filters":    appliedFilters,
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, response)

	h.Logger.PrintInfo("Categories retrieved successfully", map[string]string{
		"page":  strconv.Itoa(pagination.Page),
		"count": strconv.Itoa(len(categories)),
		"total": strconv.Itoa(totalCount),
	})
}
