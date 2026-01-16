package deletevote

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/app"
	votecommands "github.com/arnald/forum/internal/app/votes/commands"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/middleware"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Response struct {
	Message string
}

type Handler struct {
	Services app.Services
	Config   *config.ServerConfig
	Logger   logger.Logger
}

func NewHandler(services app.Services, config *config.ServerConfig, logger logger.Logger) *Handler {
	return &Handler{
		Services: services,
		Config:   config,
		Logger:   logger,
	}
}

func (h *Handler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.PrintError(logger.ErrInvalidRequestMethod, nil)
		helpers.RespondWithError(
			w,
			http.StatusMethodNotAllowed,
			"invalid method",
		)
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		h.Logger.PrintError(logger.ErrUserNotFoundInContext, nil)
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.Config.Timeouts.HandlerTimeouts.UserRegister)
	defer cancel()

	voteID, err := helpers.GetQueryInt(r, "id")
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
	}

	// val := validator.New()

	// valData := struct {
	// 	VoteID int
	// }{
	// 	VoteID: voteID,
	// }

	// validator.ValidateDeleteVote(val, valData)

	// if !val.Valid() {
	// 	h.Logger.PrintError(logger.ErrValidationFailed, val.Errors)
	// 	helpers.RespondWithError(
	// 		w,
	// 		http.StatusBadRequest,
	// 		val.ToStringErrors(),
	// 	)
	// 	return
	// }

	err = h.Services.UserServices.Commands.DeleteVote.Handle(ctx, votecommands.DeleteVoteRequest{
		VoteID: voteID,
		UserID: user.ID,
	})
	if err != nil {
		h.Logger.PrintError(err, nil)
		helpers.RespondWithError(
			w,
			http.StatusInternalServerError,
			"failed to delete vote",
		)
		return
	}

	voteResponse := Response{
		Message: "Vote deleted successfully",
	}

	helpers.RespondWithJSON(
		w,
		http.StatusOK,
		nil,
		voteResponse,
	)

	h.Logger.PrintInfo(
		"Vote deleted successfully",
		nil,
	)
}
