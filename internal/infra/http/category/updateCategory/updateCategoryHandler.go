package updatecategory

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
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          int    `json:"id"`
}

type ResponseModel struct {
	CategoryName string `json:"categoryName"`
	Message      string `json:"message"`
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

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	user := middleware.GetUserFromContext(r)

	var categoryToUpdate RequestModel

	_, err := helpers.ParseBodyRequest(r, &categoryToUpdate)
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}
	defer r.Body.Close()

	err = h.UserServices.UserServices.Commands.UpdateCategory.Handle(ctx, categorycommands.UpdateCategoryRequest{
		ID:          categoryToUpdate.ID,
		Name:        categoryToUpdate.Name,
		Description: categoryToUpdate.Description,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Error updating category",
		)
		return
	}

	response := ResponseModel{
		CategoryID: categoryToUpdate.ID,
		Message:    "Category updated successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		response)

	h.Logger.PrintInfo(
		"Category updated successfully",
		map[string]string{
			"cat_id":   strconv.Itoa(categoryToUpdate.ID),
			"cat_name": categoryToUpdate.Name,
			"user_id":  user.ID,
		})
}
