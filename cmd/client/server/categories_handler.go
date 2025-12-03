package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

// CategoriesPage handles requests to the dedicated categories page.
func (cs *ClientServer) CategoriesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := getQueryIntOr(r, "page", 1)
	search := getQueryStringOr(r, "search", "")
	orderBy := getQueryStringOr(r, "order_by", "created_at")
	order := getQueryStringOr(r, "order", "desc")
	pageSize := 12 // Categories per page

	if page < 1 {
		page = 1
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	categoriesRequest := &categoriesRequest{
		OrderBy: orderBy,
		Order:   order,
		Search:  search,
	}

	backendURL, err := createURLWithParams(backendGetCategoriesDomain, categoriesRequest)
	if err != nil {
		http.Error(w, "Error creating URL Params", http.StatusInternalServerError)
		return
	}
	// Add paginationto url
	backendURL += fmt.Sprintf("&page=%d&limit=%d", page, pageSize)

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendURL, nil)
	if err != nil {
		http.Error(w, "Error making the request", http.StatusInternalServerError)
		return
	}

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		http.Error(w, "Error with the response", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	var categoryData response
	err = DecodeBackendResponse(backendResp, &categoryData)
	if err != nil {
		http.Error(w, "Error decoding the response to json", http.StatusInternalServerError)
		return
	}

	if len(categoryData.Categories) == 0 {
		categoryData.Categories = []domain.Category{}
	} else {
		categoryData.Categories = helpers.PrepareCategories(categoryData.Categories)
	}

	user := middleware.GetUserFromContext(r.Context())

	categoryData.ActivePage = "categories"
	categoryData.User = user

	// Convert pagination values to integers for template comparison
	if paginationMap, ok := categoryData.Pagination.(map[string]interface{}); ok {
		if totalPages, ok := paginationMap["totalPages"].(float64); ok {
			paginationMap["totalPages"] = int(totalPages)
		}
		if totalItems, ok := paginationMap["totalItems"].(float64); ok {
			paginationMap["totalItems"] = int(totalItems)
		}
		if page, ok := paginationMap["page"].(float64); ok {
			paginationMap["page"] = int(page)
		}
		if nextPage, ok := paginationMap["next_page"].(float64); ok {
			paginationMap["next_page"] = int(nextPage)
		}
		if prevPage, ok := paginationMap["prev_page"].(float64); ok {
			paginationMap["prev_page"] = int(prevPage)
		}
	}

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/all_categories.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		log.Println("Error loading templates:", err)
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", categoryData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// Helper functions for query parameters.
func getQueryIntOr(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getQueryStringOr(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}
