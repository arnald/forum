package userlogin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	UserServices   app.Services
	SessionManager user.SessionManager
	Config         *config.ServerConfig
}

func NewHandler(config *config.ServerConfig, app app.Services, sm user.SessionManager) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
	}
}

type LoginUserReguestModel struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request method %v\n", r.Method)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var userToLogin LoginUserReguestModel

	err := json.NewDecoder(r.Body).Decode(&userToLogin)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"unable to decode json request",
		)
		return
	}

	defer r.Body.Close()
	user, err := h.UserServices.UserServices.Queries.UserLogin.Handle(ctx, queries.UserLoginRequest{
		Email:    userToLogin.Email,
		Password: userToLogin.Password,
	})
	if err != nil {
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Error logging in user: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "error logging in user")
		return
	}

	newSession, err := h.SessionManager.CreateSession(ctx, user.ID)
	if err != nil {
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Error creating session: %v\n", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "error creating session")
		return
	}

	cookie := h.SessionManager.NewSessionCookie(newSession.Token)
	http.SetCookie(w, cookie)

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		"User logged in successfully",
	)
}
