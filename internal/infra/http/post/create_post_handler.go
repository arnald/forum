package post

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arnald/forum/internal/app/post/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type CreatePostHandler struct {
	createPostHandler queries.CreatePostRequestHandler
}

func NewCreatePostHandler(createPostHandler queries.CreatePostRequestHandler) *CreatePostHandler {
	return &CreatePostHandler{
		createPostHandler: createPostHandler,
	}
}

func (h *CreatePostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
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

	createPostReq := queries.CreatePostRequest{
		Title:      strings.TrimSpace(req.Title),
		Content:    strings.TrimSpace(req.Content),
		UserID:     userID,
		Categories: req.Categories,
	}

	post, err := h.createPostHandler.Handle(r.Context(), createPostReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyTitle:
			helpers.RespondWithError(w, http.StatusBadRequest, "Title is required")
		case queries.ErrEmptyContent:
			helpers.RespondWithError(w, http.StatusBadRequest, "Content is required")
		case queries.ErrEmptyUserID:
			helpers.RespondWithError(w, http.StatusUnauthorized, "User authentication required")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, nil, map[string]interface{}{
		"success": true,
		"post":    post,
	})
}