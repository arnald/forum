package server

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

const defaultPageSize = 10

type topicsRequest struct {
	OrderBy  string `url:"order_by"`
	Order    string `url:"order"`
	Search   string `url:"search"`
	Category int    `url:"category"`
	Page     int    `url:"page"`
	PageSize int    `url:"page_size"`
}

type topicsResponse struct {
	User       *domain.LoggedInUser   `json:"user"`
	Filters    map[string]interface{} `json:"filters"`
	Topics     []domain.Topic         `json:"topics"`
	Categories []domain.Category      `json:"categories"`
	Pagination domain.Pagination      `json:"pagination"`
}

// TopicsPage handles GET requests to /topics.
func (cs *ClientServer) TopicsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := getQueryIntOr(r, "page", 1)
	search := getQueryStringOr(r, "search", "")
	orderBy := getQueryStringOr(r, "order_by", "created_at")
	order := getQueryStringOr(r, "order", "desc")
	category := getQueryIntOr(r, "category", 0)
	pageSize := getQueryIntOr(r, "page_size", defaultPageSize)

	topicsReq := &topicsRequest{
		OrderBy:  orderBy,
		Order:    order,
		Search:   search,
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}

	backendURL, err := createURLWithParams(backendGetTopicsDomain(), topicsReq)
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

	var pageData topicsResponse
	err = helpers.DecodeBackendResponse(backendResp, &pageData)
	if err != nil {
		http.Error(w, "Error with decoding response into data struct", http.StatusInternalServerError)
		return
	}

	if len(pageData.Topics) == 0 {
		pageData.Topics = []domain.Topic{}
	}

	for i := range pageData.Topics {
		// Normalize all category colors in the slice
		for j := range pageData.Topics[i].CategoryColors {
			pageData.Topics[i].CategoryColors[j] = helpers.NormalizeColor(pageData.Topics[i].CategoryColors[j])
		}

		// // Also normalize the single CategoryColor for backward compatibility (if you're still using it)
		// if pageData.Topics[i].CategoryColor != "" {
		// 	pageData.Topics[i].CategoryColor = helpers.NormalizeColor(pageData.Topics[i].CategoryColor)
		// }
	}
	pageData.User = middleware.GetUserFromContext(r.Context())

	// Create template with custom functions
	tmpl := template.New("base").Funcs(template.FuncMap{
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
	})

	// Parse files with the custom functions available
	tmpl, err = tmpl.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/all_topics.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
