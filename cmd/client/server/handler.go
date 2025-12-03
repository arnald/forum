package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"text/template"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

type categoriesRequest struct {
	OrderBy string
	Order   string
	Search  string
}

type response struct {
	Filters    any                  `json:"filters"`
	Pagination any                  `json:"pagination"`
	User       *domain.LoggedInUser `json:"user"`
	ActivePage string               `json:"activePage"`
	Categories []domain.Category    `json:"categories"`
}

// backendError is a custom error type for backend errors.
type backendError string

func (e backendError) Error() string {
	return string(e)
}

const backendGetCategoriesDomain = "http://localhost:8080/api/v1/categories/all"

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

	backendURL, err := createURLWithParams(backendGetCategoriesDomain, defaultCategoriesOptions)
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
	defer backendResp.Body.Close()

	var categoryData response
	err = DecodeBackendResponse(backendResp, &categoryData)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	categoryData.Categories = helpers.PrepareCategories(categoryData.Categories)

	user := middleware.GetUserFromContext(r.Context())

	categoryData.ActivePage = "home"
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
		log.Println("Error loading home.html:", err)
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", categoryData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
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

		fieldName := strings.ToLower(fieldType.Name[:1]) + fieldType.Name[1:]
		fieldValue := fmt.Sprintf("%v", field.Interface())

		queryParams.Add(fieldName, fieldValue)
	}

	completeURL := domainURL + "?" + queryParams.Encode()

	return completeURL, nil
}
