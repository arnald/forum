package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/infra/session"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	UserServices   app.Services
	SessionManager *session.SessionManager
}

func NewHandler(app app.Services, sm *session.SessionManager) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
	}
}

type RegisterUserReguestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request method %v\n", r.Method)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var userToRegister RegisterUserReguestModel

	err := json.NewDecoder(r.Body).Decode(&userToRegister)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"unable to decode json request",
		)
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Hour*3))
	defer cancel()

	user, err := h.UserServices.UserServices.Queries.UserRegister.Handle(ctx, queries.UserRegisterRequest{
		Name:     userToRegister.Username,
		Password: userToRegister.Password,
		Email:    userToRegister.Email,
	})
	if err != nil {
		switch {
		case errors.Is(err, sqlite.ErrDuplicateEmail):
			helpers.RespondWithError(
				w,
				http.StatusConflict,
				"a user with this email address already exists",
			)
		default:
			helpers.RespondWithError(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
		}
		return
	}

	newSession, err := h.SessionManager.CreateSession(user.ID)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	helpers.RespondWithJSON(
		w,
		http.StatusCreated,
		nil,
		newSession,
	)
}
