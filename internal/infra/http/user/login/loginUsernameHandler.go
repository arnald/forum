//nolint:dupl
package userlogin

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type LoginUserUsernameRequestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h Handler) UserLoginUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	//TODO : FIND WHERE TO IMPLEMENT A COMMON KEY FOR USER CONTEXT
	if _, err := helpers.GetUserFromContext(r.Context(), "user"); err == nil {
		h.Logger.PrintInfo("Request made by authenticated user", nil)
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	var userToLogin LoginUserUsernameRequestModel

	userAny, err := helpers.ParseBodyRequest(r, &userToLogin)
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

	validator.ValidateUserLoginUsername(v, userAny)

	if !v.Valid() {
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			v.ToStringErrors(),
		)

		h.Logger.PrintError(logger.ErrValidationFailed, v.Errors)

		return
	}

	user, err := h.UserServices.UserServices.Queries.UserLoginUsername.Handle(ctx, queries.UserLoginUsernameRequest{
		Username: userToLogin.Username,
		Password: userToLogin.Password,
	})
	if err != nil {
		helpers.RespondWithError(w,
			http.StatusInternalServerError,
			"error logging in user",
		)

		h.Logger.PrintError(err, nil)
		return
	}

	newSession, err := h.SessionManager.CreateSession(ctx, user.ID)
	if err != nil {
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"error creating session",
		)

		h.Logger.PrintError(err, nil)
		return
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		newSession,
	)
}
