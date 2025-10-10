package createcategory

import (
	"context"
	"net/http"

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
}

type ResponseModel struct {
	CategoryName string `json:"categoryName"`
	Message      string `json:"message"`
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

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	user := middleware.GetUserFromContext(r)

	var categoryToCreate RequestModel

	_, err := helpers.ParseBodyRequest(r, &categoryToCreate)
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}
	defer r.Body.Close()

	err = h.UserServices.UserServices.Commands.CreateCategory.Handle(ctx, categorycommands.CreateCategoryRequest{
		Name:        categoryToCreate.Name,
		Description: categoryToCreate.Description,
		CreatedBy:   user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"Failed to create category",
		)

		h.Logger.PrintError(err, nil)
		return
	}

	response := ResponseModel{
		CategoryName: categoryToCreate.Name,
		Message:      "Category created successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		response,
	)

	h.Logger.PrintInfo(
		"Category created successfully",
		map[string]string{
			"cat_name": categoryToCreate.Name,
			"user_id":  user.ID,
		})
}
