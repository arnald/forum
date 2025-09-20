package comment

import (
	"net/http"

	"github.com/arnald/forum/internal/app/comment/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type GetCommentsHandler struct {
	getCommentsByPostHandler queries.GetCommentsByPostRequestHandler
	getCommentTreeHandler    queries.GetCommentTreeRequestHandler
}

func NewGetCommentsHandler(
	getCommentsByPostHandler queries.GetCommentsByPostRequestHandler,
	getCommentTreeHandler queries.GetCommentTreeRequestHandler,
) *GetCommentsHandler {
	return &GetCommentsHandler{
		getCommentsByPostHandler: getCommentsByPostHandler,
		getCommentTreeHandler:    getCommentTreeHandler,
	}
}

func (h *GetCommentsHandler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	comments, err := h.getCommentsByPostHandler.Handle(r.Context(), queries.GetCommentsByPostRequest{PostID: postID})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":  true,
		"comments": comments,
	})
}

func (h *GetCommentsHandler) GetCommentTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	commentTree, err := h.getCommentTreeHandler.Handle(r.Context(), queries.GetCommentTreeRequest{PostID: postID})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comment tree")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":     true,
		"commentTree": commentTree,
	})
}