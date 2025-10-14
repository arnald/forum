package deletecategory

import (
	"context"
	"net/http"
	"strconv"

	"github.com/arnald/forum/internal/app"
	categorycommands "github.com/arnald/forum/internal/app/categories/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type RequestModel struct {
	CategoryID int `json:"categoryId"`
}

type ResponseModel struct {
	Message    string `json:"message"`
	CategoryID int    `json:"categoryId"`
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

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	user := middleware.GetUserFromContext(r)

	var categoryToDelete RequestModel

	_, err := helpers.ParseBodyRequest(r, &categoryToDelete)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.UserServices.UserServices.Commands.DeleteCategory.Handle(ctx, categorycommands.DeleteCategoryRequest{
		CategoryID: categoryToDelete.CategoryID,
		UserID:     user.ID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting category")
		return
	}

	helpers.RespondWithJSON(w,
		http.StatusOK,
		nil,
		ResponseModel{
			CategoryID: categoryToDelete.CategoryID,
			Message:    "Category deleted successfully",
		})

	h.Logger.PrintInfo(
		"Category deleted successfully",
		map[string]string{
			"cat_id":  strconv.Itoa(categoryToDelete.CategoryID),
			"user_id": user.ID,
		})
}
