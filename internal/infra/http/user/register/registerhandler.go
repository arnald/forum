package userregister

import (
	"context"
	"net/http"
	"strings"

	"github.com/arnald/forum/internal/app"
	userQueries "github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type RegisterUserReguestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterUserResponse struct {
	UserID  string `json:"userdId"`
	Message string `json:"message"`
}

type Handler struct {
	UserServices   app.Services
	SessionManager session.Manager
	Config         *config.ServerConfig
	Logger         logger.Logger
}

func NewHandler(config *config.ServerConfig, app app.Services, sm session.Manager, logger logger.Logger) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
		Logger:         logger,
	}
}

func (h Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var userToRegister RegisterUserReguestModel

	userAny, err := helpers.ParseBodyRequest(r, &userToRegister)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid request: "+err.Error(),
		)

		h.Logger.PrintError(err, nil)

		return
	}
	defer r.Body.Close()

	v := validator.New()

	validator.ValidateUserRegistration(v, userAny)

	if !v.Valid() {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			v.ToStringErrors(),
		)

		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)
		return
	}

	user, err := h.UserServices.UserServices.Queries.UserRegister.Handle(ctx, userQueries.UserRegisterRequest{
		Name:     userToRegister.Username,
		Password: userToRegister.Password,
		Email:    strings.ToLower(userToRegister.Email),
	})
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)

		h.Logger.PrintError(err, nil)

		return
	}

	userResponse := RegisterUserResponse{
		UserID:  user.ID,
		Message: "user registered successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		userResponse,
	)

	h.Logger.PrintInfo(
		"User registered successfully",
		map[string]string{
			"userId": user.ID,
			"email":  user.Email,
			"name":   user.Username,
		},
	)
}
