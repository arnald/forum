package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

const minURLPathLength = 2

type topicPageResponse struct {
	UserVote       *int             `json:"userVote"`
	ImagePath      string           `json:"imagePath"`
	OwnerUsername  string           `json:"ownerUsername"`
	Content        string           `json:"content"`
	UserID         string           `json:"userId"`
	CreatedAt      string           `json:"createdAt"`
	Title          string           `json:"title"`
	UpdatedAt      string           `json:"updatedAt"`
	CategoryColors []string         `json:"categoryColors"`
	CategoryNames  []string         `json:"categoryNames"`
	Comments       []domain.Comment `json:"comments"`
	CategoryIDs    []int            `json:"categoryIds"`
	Upvotes        int              `json:"upvotes"`
	Downvotes      int              `json:"downvotes"`
	Score          int              `json:"score"`
	TopicID        int              `json:"topicId"`
}

type topicPageRequest struct {
	TopicID string `url:"id"`
}

type topicPageData struct {
	User       *domain.LoggedInUser `json:"user"`
	Categories []domain.Category    `json:"categories"`
	Topic      domain.Topic         `json:"topic"`
}

// TopicPage handles GET requests to /topic/{id}.
func (cs *ClientServer) TopicPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < minURLPathLength {
		templates.NotFoundHandler(w, r, "Invalid topic URL", http.StatusBadRequest)
		return
	}

	topicIDStr := pathParts[1]
	_, err := strconv.Atoi(topicIDStr)
	if err != nil {
		templates.NotFoundHandler(w, r, "Invalid topic ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	topicReq := &topicPageRequest{
		TopicID: topicIDStr,
	}

	topicURL, err := createURLWithParams(backendGetTopicByID, topicReq)
	if err != nil {
		log.Printf("Error creating topic URL: %v", err)
		http.Error(w, "Error creating URL params", http.StatusInternalServerError)
		return
	}

	topicHTTPReq, err := http.NewRequestWithContext(ctx, http.MethodGet, topicURL, nil)
	if err != nil {
		log.Printf("Error creating topic request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		topicHTTPReq.AddCookie(cookie)
	}

	topicResp, err := cs.HTTPClient.Do(topicHTTPReq)
	if err != nil {
		log.Printf("Error fetching topic: %v", err)
		http.Error(w, "Error with the response", http.StatusInternalServerError)
		return
	}
	defer topicResp.Body.Close()

	if topicResp.StatusCode == http.StatusNotFound {
		templates.NotFoundHandler(w, r, "Topic not found", http.StatusNotFound)
		return
	}

	if topicResp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %d", topicResp.StatusCode)
		templates.NotFoundHandler(w, r, "Error loading topic", http.StatusInternalServerError)
		return
	}

	var topicData topicPageResponse
	err = helpers.DecodeBackendResponse(topicResp, &topicData)
	if err != nil {
		log.Printf("Error decoding topic response: %v", err)
		http.Error(w, "Error with decoding response into data struct", http.StatusInternalServerError)
		return
	}

	// Fetch categories for the edit form
	categoriesHTTPReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendGetCategoriesDomain, nil)
	if err != nil {
		log.Printf("Error creating categories request: %v", err)
		http.Error(w, "Error creating categories request", http.StatusInternalServerError)
		return
	}

	categoriesResp, err := cs.HTTPClient.Do(categoriesHTTPReq)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}
	defer categoriesResp.Body.Close()

	var categoriesData struct {
		Categories []domain.Category `json:"categories"`
	}

	if categoriesResp.StatusCode == http.StatusOK {
		err = helpers.DecodeBackendResponse(categoriesResp, &categoriesData)
		if err != nil {
			log.Printf("Error decoding categories response: %v", err)
			// Continue without categories - edit form won't work but page will still load
			categoriesData.Categories = []domain.Category{}
		}
	} else {
		log.Printf("Failed to fetch categories, status: %d", categoriesResp.StatusCode)
		categoriesData.Categories = []domain.Category{}
	}

	normalizedColors := make([]string, len(topicData.CategoryColors))
	for i, color := range topicData.CategoryColors {
		normalizedColors[i] = helpers.NormalizeColor(color)
	}

	topic := domain.Topic{
		ID:             topicData.TopicID,
		CategoryIDs:    topicData.CategoryIDs,
		Title:          topicData.Title,
		Content:        topicData.Content,
		ImagePath:      topicData.ImagePath,
		UserID:         topicData.UserID,
		CreatedAt:      topicData.CreatedAt,
		UpdatedAt:      topicData.UpdatedAt,
		UpvoteCount:    topicData.Upvotes,
		DownvoteCount:  topicData.Downvotes,
		VoteScore:      topicData.Score,
		UserVote:       topicData.UserVote,
		OwnerUsername:  topicData.OwnerUsername,
		Comments:       topicData.Comments,
		CategoryNames:  topicData.CategoryNames,
		CategoryColors: normalizedColors,
	}

	pageData := topicPageData{
		User:       middleware.GetUserFromContext(r.Context()),
		Topic:      topic,
		Categories: categoriesData.Categories,
	}

	tmpl, err := template.New("base").
		Funcs(template.FuncMap{
			"hasID": hasID,
		}).
		ParseFiles(
			"frontend/html/layouts/base.html",
			"frontend/html/pages/topic.html",
			"frontend/html/partials/navbar.html",
			"frontend/html/partials/footer.html",
		)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func hasID(ids []int, id int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}
