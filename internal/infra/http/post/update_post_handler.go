package post

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arnald/forum/internal/app/post/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type UpdatePostHandler struct {
	updatePostHandler queries.UpdatePostRequestHandler
}

func NewUpdatePostHandler(updatePostHandler queries.UpdatePostRequestHandler) *UpdatePostHandler {
	return &UpdatePostHandler{
		updatePostHandler: updatePostHandler,
	}
}

func (h *UpdatePostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ID         string   `json:"id"`
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
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

	updatePostReq := queries.UpdatePostRequest{
		ID:         req.ID,
		Title:      strings.TrimSpace(req.Title),
		Content:    strings.TrimSpace(req.Content),
		UserID:     userID,
		Categories: req.Categories,
	}

	post, err := h.updatePostHandler.Handle(r.Context(), updatePostReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyTitle:
			helpers.RespondWithError(w, http.StatusBadRequest, "Title is required")
		case queries.ErrEmptyContent:
			helpers.RespondWithError(w, http.StatusBadRequest, "Content is required")
		case queries.ErrPostNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, "Post not found or you don't have permission to edit it")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to update post")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success": true,
		"post":    post,
	})
}