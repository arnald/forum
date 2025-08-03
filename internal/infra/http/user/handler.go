package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/session"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	UserServices   app.Services
	SessionManager *session.Manager
	Config         *config.ServerConfig
}

func NewHandler(config *config.ServerConfig, app app.Services, sm *session.Manager) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
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
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request method %v\n", r.Method)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	// decBody := json.NewDecoder(r.Body)
	// decBody.DisallowUnknownFields()
	var userToRegister RegisterUserReguestModel

	// err := decBody.Decode(&userToRegister)
	userAny, err := helpers.ParseBodyRequest(r, &userToRegister)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("invalid request: %s", err.Error()),
		)

		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request:  %v\n", err.Error())

		return
	}
	defer r.Body.Close()

	userRegister, ok := userAny.(*RegisterUserReguestModel)
	if !ok {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			"email, username and password are required",
		)
		return
	}

	err = helpers.UserRegisterIsValid(userRegister.Username, userRegister.Password, userRegister.Email)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	user, err := h.UserServices.UserServices.Queries.UserRegister.Handle(ctx, queries.UserRegisterRequest{
		Name:     userToRegister.Username,
		Password: userToRegister.Password,
		Email:    userToRegister.Email,
	})
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	newSession, err := h.SessionManager.CreateSession(ctx, user.ID)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	sessionResponse := &RegisterUserSessionResponse{
		AccessToken:  newSession.AccessToken,
		RefreshToken: newSession.RefreshToken,
		UserID:       newSession.UserID,
	}
	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		sessionResponse,
	)
}
