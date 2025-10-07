package userlogin

import (
	"context"
	"net/http"

	userQueries "github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/pkg/helpers"
	"github.com/arnald/forum/internal/pkg/validator"
)

type LoginUserUsernameRequestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h Handler) UserLoginUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
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

	user, err := h.UserServices.UserServices.Queries.UserLoginUsername.Handle(ctx, userQueries.UserLoginUsernameRequest{
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

	loginResponse := LoginResponse{
		UserID:       user.ID,
		Username:     user.Username,
		AccessToken:  newSession.AccessToken,
		RefreshToken: newSession.RefreshToken,
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		loginResponse,
	)

	h.Logger.PrintInfo(
		"User login successfully",
		map[string]string{
			"userId": user.ID,
			"name":   user.Username,
		},
	)
}
