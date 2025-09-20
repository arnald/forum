package comment

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arnald/forum/internal/app/comment/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type UpdateCommentHandler struct {
	updateCommentHandler queries.UpdateCommentRequestHandler
}

func NewUpdateCommentHandler(updateCommentHandler queries.UpdateCommentRequestHandler) *UpdateCommentHandler {
	return &UpdateCommentHandler{
		updateCommentHandler: updateCommentHandler,
	}
}

func (h *UpdateCommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	updateCommentReq := queries.UpdateCommentRequest{
		ID:      req.ID,
		Content: strings.TrimSpace(req.Content),
		UserID:  userID,
	}

	comment, err := h.updateCommentHandler.Handle(r.Context(), updateCommentReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyContent:
			helpers.RespondWithError(w, http.StatusBadRequest, "Comment content is required")
		case queries.ErrCommentNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, "Comment not found or you don't have permission to edit it")
		case queries.ErrEmptyUserID:
			helpers.RespondWithError(w, http.StatusUnauthorized, "User authentication required")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to update comment")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success": true,
		"comment": comment,
	})
}