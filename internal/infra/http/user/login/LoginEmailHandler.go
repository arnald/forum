//nolint:dupl
package userlogin

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type LoginUserEmailRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h Handler) UserLoginEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
		logger.Printf("Invalid request method %v\n", r.Method)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var userToLogin LoginUserEmailRequestModel

	userAny, err := helpers.ParseBodyRequest(r, &userToLogin)
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

	validator.ValidateUserLoginEmail(v, userAny)

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

	user, err := h.UserServices.UserServices.Queries.UserLoginEmail.Handle(ctx, queries.UserLoginEmailRequest{
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

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		newSession,
	)
}
