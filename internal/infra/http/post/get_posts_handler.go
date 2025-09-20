package post

import (
	"net/http"

	"github.com/arnald/forum/internal/app/post/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type GetPostsHandler struct {
	getAllPostsHandler      queries.GetAllPostsRequestHandler
	getPostByIDHandler      queries.GetPostByIDRequestHandler
	getPostsByCategoryHandler queries.GetPostsByCategoryRequestHandler
}

func NewGetPostsHandler(
	getAllPostsHandler queries.GetAllPostsRequestHandler,
	getPostByIDHandler queries.GetPostByIDRequestHandler,
	getPostsByCategoryHandler queries.GetPostsByCategoryRequestHandler,
) *GetPostsHandler {
	return &GetPostsHandler{
		getAllPostsHandler:        getAllPostsHandler,
		getPostByIDHandler:        getPostByIDHandler,
		getPostsByCategoryHandler: getPostsByCategoryHandler,
	}
}

func (h *GetPostsHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Check if category filtering is requested
	categoryID := r.URL.Query().Get("category_id")
	categoryName := r.URL.Query().Get("category")

	var posts interface{}
	var err error

	if categoryID != "" || categoryName != "" {
		// Get posts by category
		posts, err = h.getPostsByCategoryHandler.Handle(r.Context(), queries.GetPostsByCategoryRequest{
			CategoryID:   categoryID,
			CategoryName: categoryName,
		})
	} else {
		// Get all posts
		posts, err = h.getAllPostsHandler.Handle(r.Context(), queries.GetAllPostsRequest{})
	}

	if err != nil {
		switch err {
		case queries.ErrEmptyCategory:
			helpers.RespondWithError(w, http.StatusBadRequest, "Category parameter cannot be empty")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch posts")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success": true,
		"posts":   posts,
	})
}

func (h *GetPostsHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	post, err := h.getPostByIDHandler.Handle(r.Context(), queries.GetPostByIDRequest{ID: postID})
	if err != nil {
		if err == queries.ErrPostNotFound {
			helpers.RespondWithError(w, http.StatusNotFound, "Post not found")
		} else {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success": true,
		"post":    post,
	})
}