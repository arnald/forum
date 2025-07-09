package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Handler struct {
	UserServices app.Services
}

func NewHandler(app app.Services) *Handler {
	return &Handler{
		UserServices: app,
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

	err = h.UserServices.UserServices.Queries.UserRegister.Handle(queries.UserRegisterRequest{
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

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		"user registered succesfully",
	)
}
