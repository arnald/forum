package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/session"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type RegisterUserResponse struct {
	UserID  string `json:"userdId"`
	Message string `json:"message"`
}

type Handler struct {
	UserServices   app.Services
	SessionManager *session.Manager
	Config         *config.ServerConfig
	Logger         logger.Logger
}

func NewHandler(config *config.ServerConfig, app app.Services, sm *session.Manager, logger logger.Logger) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
		Logger:         logger,
	}
}

type RegisterUserReguestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterUserSessionResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
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

		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request:  %v\n", err.Error())

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

		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Println("Invalid request: " + v.ToStringErrors())

		return
	}

	user, err := h.UserServices.UserServices.Queries.UserRegister.Handle(ctx, queries.UserRegisterRequest{
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

		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Println(err.Error())
		return
	}

	userRegistered := RegisterUserResponse{
		UserID:  user.ID,
		Message: "Account was successfully created!",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		userRegistered,
	)
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logger.Println("user with id: " + user.ID + "and username: " + userToRegister.Username + " and email " + userToRegister.Email + "was successfully created!")
}
