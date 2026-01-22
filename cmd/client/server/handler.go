package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

type categoriesRequest struct {
	OrderBy  string `url:"order_by"`
	Order    string `url:"order"`
	Search   string `url:"search"`
	Page     int    `url:"page"`
	PageSize int    `url:"page_size"`
}

type response struct {
	Filters    any                  `json:"filters"`
	User       *domain.LoggedInUser `json:"user"`
	Categories []domain.Category    `json:"categories"`
	Pagination domain.Pagination    `json:"pagination"`
}

// backendError is a custom error type for backend errors.
type backendError string

func (e backendError) Error() string {
	return string(e)
}

// HomePage handles requests to the homepage.
func (cs *ClientServer) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defaultCategoriesOptions := &categoriesRequest{
		OrderBy: "created_at",
		Order:   "desc",
		Search:  "",
	}

	backendURL, err := createURLWithParams(backendGetCategoriesDomain(), defaultCategoriesOptions)
	if err != nil {
		http.Error(w, "Error creating URL params", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendURL, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(httpReq, ip)

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		http.Error(w, "Error with the response", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	var categoryData response
	err = helpers.DecodeBackendResponse(backendResp, &categoryData)
	if err != nil {
		log.Printf("Decode error: %v, Status: %d, URL: %s", err, backendResp.StatusCode, backendURL)
		http.Error(w, "Error with decoding response into data struct", http.StatusInternalServerError)
		return
	}

	categoryData.Categories = helpers.PrepareCategories(categoryData.Categories)

	user := middleware.GetUserFromContext(r.Context())

	categoryData.User = user

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/home.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/category_details.html",
		"frontend/html/partials/categories.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", categoryData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

var ErrFailedToCreateURL = errors.New("failed to create url with params")

func createURLWithParams(domainURL string, params any) (string, error) {
	val := reflect.ValueOf(params).Elem()
	if !val.IsValid() {
		return "", ErrFailedToCreateURL
	}

	typ := val.Type()
	queryParams := url.Values{}

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)

		fieldName := fieldType.Tag.Get("url")
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name[:1]) + fieldType.Name[1:]
			fieldName = toSnakeCase(fieldName)
		}

		fieldValue := fmt.Sprintf("%v", field.Interface())

		queryParams.Add(fieldName, fieldValue)
	}

	completeURL := domainURL + "?" + queryParams.Encode()
	return completeURL, nil
}

func toSnakeCase(input string) string {
	var result strings.Builder
	for i, char := range input {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(char))
	}
	return result.String()
}
