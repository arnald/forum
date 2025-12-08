package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

const defaltPageSize = 12

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
	pageSize := getQueryIntOr(r, "page_size", defaltPageSize)

	categoriesRequest := &categoriesRequest{
		OrderBy:  orderBy,
		Order:    order,
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	backendURL, err := createURLWithParams(backendGetCategoriesDomain, categoriesRequest)
	if err != nil {
		http.Error(w, "Error creating URL Params", http.StatusInternalServerError)
		return
	}

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
	err = helpers.DecodeBackendResponse(backendResp, &categoryData)
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

	categoryData.User = user

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
