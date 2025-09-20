package category

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/arnald/forum/internal/app/category/queries"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type CategoryHandler struct {
	createCategoryHandler         queries.CreateCategoryRequestHandler
	getAllCategoriesHandler       queries.GetAllCategoriesRequestHandler
	getCategoryByIDHandler        queries.GetCategoryByIDRequestHandler
	updateCategoryHandler         queries.UpdateCategoryRequestHandler
	deleteCategoryHandler         queries.DeleteCategoryRequestHandler
	getCategoryWithPostsHandler   queries.GetCategoryWithPostsRequestHandler
}

func NewCategoryHandler(
	createCategoryHandler queries.CreateCategoryRequestHandler,
	getAllCategoriesHandler queries.GetAllCategoriesRequestHandler,
	getCategoryByIDHandler queries.GetCategoryByIDRequestHandler,
	updateCategoryHandler queries.UpdateCategoryRequestHandler,
	deleteCategoryHandler queries.DeleteCategoryRequestHandler,
	getCategoryWithPostsHandler queries.GetCategoryWithPostsRequestHandler,
) *CategoryHandler {
	return &CategoryHandler{
		createCategoryHandler:       createCategoryHandler,
		getAllCategoriesHandler:     getAllCategoriesHandler,
		getCategoryByIDHandler:      getCategoryByIDHandler,
		updateCategoryHandler:       updateCategoryHandler,
		deleteCategoryHandler:       deleteCategoryHandler,
		getCategoryWithPostsHandler: getCategoryWithPostsHandler,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	createCategoryReq := queries.CreateCategoryRequest{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	}

	category, err := h.createCategoryHandler.Handle(r.Context(), createCategoryReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyName:
			helpers.RespondWithError(w, http.StatusBadRequest, "Category name is required")
		case queries.ErrCategoryAlreadyExists:
			helpers.RespondWithError(w, http.StatusConflict, "Category with this name already exists")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create category")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, nil, map[string]interface{}{
		"success":  true,
		"category": category,
	})
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	categories, err := h.getAllCategoriesHandler.Handle(r.Context(), queries.GetAllCategoriesRequest{})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":    true,
		"categories": categories,
	})
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	category, err := h.getCategoryByIDHandler.Handle(r.Context(), queries.GetCategoryByIDRequest{ID: categoryID})
	if err != nil {
		switch err {
		case queries.ErrCategoryNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, "Category not found")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch category")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":  true,
		"category": category,
	})
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updateCategoryReq := queries.UpdateCategoryRequest{
		ID:          req.ID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	}

	category, err := h.updateCategoryHandler.Handle(r.Context(), updateCategoryReq)
	if err != nil {
		switch err {
		case queries.ErrEmptyID:
			helpers.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		case queries.ErrEmptyName:
			helpers.RespondWithError(w, http.StatusBadRequest, "Category name is required")
		case queries.ErrCategoryNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, "Category not found")
		case queries.ErrCategoryAlreadyExists:
			helpers.RespondWithError(w, http.StatusConflict, "Category with this name already exists")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to update category")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":  true,
		"category": category,
	})
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	err := h.deleteCategoryHandler.Handle(r.Context(), queries.DeleteCategoryRequest{ID: categoryID})
	if err != nil {
		switch err {
		case queries.ErrCategoryNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, "Category not found")
		case queries.ErrCategoryHasPosts:
			helpers.RespondWithError(w, http.StatusConflict, "Cannot delete category that contains posts")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to delete category")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success": true,
		"message": "Category deleted successfully",
	})
}

func (h *CategoryHandler) GetCategoryWithPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helpers.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	categoryID := r.URL.Query().Get("id")
	if categoryID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // default offset
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	categoryWithPosts, err := h.getCategoryWithPostsHandler.Handle(r.Context(), queries.GetCategoryWithPostsRequest{
		CategoryID: categoryID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		switch err {
		case queries.ErrCategoryNotFound:
			helpers.RespondWithError(w, http.StatusNotFound, "Category not found")
		default:
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch category with posts")
		}
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, nil, map[string]interface{}{
		"success":  true,
		"category": categoryWithPosts.Category,
		"posts":    categoryWithPosts.Posts,
	})
}