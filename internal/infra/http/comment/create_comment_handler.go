package comment

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arnald/forum/internal/app/comment/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type CreateCommentHandler struct {
	createCommentHandler queries.CreateCommentRequestHandler
}

func NewCreateCommentHandler(createCommentHandler queries.CreateCommentRequestHandler) *CreateCommentHandler {
	return &CreateCommentHandler{
		createCommentHandler: createCommentHandler,
	}
}

func (h *CreateCommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Content  string  `json:"content"`
		PostID   string  `json:"post_id"`
		ParentID *string `json:"parent_id,omitempty"`
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

	createCommentReq := queries.CreateCommentRequest{
		Content:  strings.TrimSpace(req.Content),
		PostID:   req.PostID,
		UserID:   userID,
		ParentID: req.ParentID,
	}

	comment, err := h.createCommentHandler.Handle(r.Context(), createCommentReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyContent:
			helpers.RespondWithError(w, http.StatusBadRequest, "Comment content is required")
		case queries.ErrEmptyPostID:
			helpers.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
		case queries.ErrEmptyUserID:
			helpers.RespondWithError(w, http.StatusUnauthorized, "User authentication required")
		case queries.ErrParentCommentNotFound:
			helpers.RespondWithError(w, http.StatusBadRequest, "Parent comment not found")
		case queries.ErrMaxNestingExceeded:
			helpers.RespondWithError(w, http.StatusBadRequest, "Maximum comment nesting level exceeded")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, nil, map[string]interface{}{
		"success": true,
		"comment": comment,
	})
}