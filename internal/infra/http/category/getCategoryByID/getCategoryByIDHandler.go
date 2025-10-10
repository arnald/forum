package getcategorybyid

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	categoryqueries "github.com/arnald/forum/internal/app/categories/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	CategoryID int `json:"categoryId"`
}

type ResponseModel struct {
	CategoryName string `json:"categoryName"`
	CategoryID   int    `json:"categoryId"`
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

func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var categoryToGet RequestModel

	_, err := helpers.ParseBodyRequest(r, &categoryToGet)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	category, err := h.UserServices.UserServices.Queries.GetCategoryByID.Handle(ctx, categoryqueries.GetCategoryByIDRequest{
		ID: categoryToGet.CategoryID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error getting category")
		return
	}

	helpers.RespondWithJSON(w,
		http.StatusOK,
		nil,
		ResponseModel{
			CategoryID:   category.ID,
			CategoryName: category.Name,
		},
	)

	h.Logger.PrintInfo(
		"Category retrieved successfully",
		map[string]string{
			"category_id": strconv.Itoa(category.ID),
			"name":        category.Name,
		},
	)
}
