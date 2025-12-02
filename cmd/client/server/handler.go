package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/middleware"
)

type categoriesRequest struct {
	Order_by string
	Order    string
	Search   string
}

type response struct {
	Filters    any `json:"filters"`
	Categories any `json:"categories"`
	Pagination any `json:"pagination"`
}

// backendError is a custom error type for backend errors.
type backendError string

func (e backendError) Error() string {
	return string(e)
}

var (
	backendGetCategoriesDomain = "http://localhost:8080/api/v1/categories/all"
)

// HomePage handles requests to the homepage.
func (cs *ClientServer) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defaultCategoriesOptions := &categoriesRequest{
		Order_by: "created_at",
		Order:    "desc",
		Search:   "",
	}

	backendURL, err := createUrlWithParams(backendGetCategoriesDomain, defaultCategoriesOptions)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendURL, nil)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var categoryData response
	err = DecodeBackendResponse(backendResp, &categoryData)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user := middleware.GetUserFromContext(r.Context())

	pageData := struct {
		User       *domain.LoggedInUser
		ActivePage string
		Categories any
		Pagination any
		Filters    any
	}{
		User:       user,
		ActivePage: "home",
		Categories: categoryData.Categories,
		Pagination: categoryData.Pagination,
		Filters:    categoryData.Filters,
	}

	// tmpl, err := template.ParseFiles(
	// 	"frontend/html/layouts/base.html",
	// 	"frontend/html/pages/home.html",
	// 	"frontend/html/partials/navbar.html",
	// 	"frontend/html/partials/category_details.html",
	// 	"frontend/html/partials/categories.html",
	// 	"frontend/html/partials/footer.html",
	// )
	// if err != nil {
	// 	log.Println("Error loading home.html:", err)
	// 	templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
	// 	return
	// }

	// err = tmpl.ExecuteTemplate(w, "base", pageData)
	// if err != nil {
	// 	log.Println("Error executing template:", err)
	// 	http.Error(w, "Failed to render page", http.StatusInternalServerError)
	// }

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pageData)
	if err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

func DecodeBackendResponse[T any](resp *http.Response, target *T) error {
	wrapper := struct {
		Data T `json:"data"`
	}{}

	err := json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return err
	}

	*target = wrapper.Data

	return nil
}

func createUrlWithParams(domainURL string, params any) (string, error) {
	val := reflect.ValueOf(params).Elem()
	if !val.IsValid() {
		return "", fmt.Errorf("failed to create url")
	}

	typ := val.Type()

	queryParams := url.Values{}
	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)

		fieldName := strings.ToLower(fieldType.Name[:1]) + fieldType.Name[1:]
		fieldValue := fmt.Sprintf("%v", field.Interface())

		queryParams.Add(fieldName, fieldValue)

	}

	completeURL := domainURL + "?" + queryParams.Encode()

	return completeURL, nil
}
