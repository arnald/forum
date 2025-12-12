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
	UserVote      *int             `json:"userVote"`
	ImagePath     string           `json:"imagePath"`
	OwnerUsername string           `json:"ownerUsername"`
	Content       string           `json:"content"`
	CategoryColor string           `json:"categoryColor"`
	UserID        string           `json:"userId"`
	CreatedAt     string           `json:"createdAt"`
	Title         string           `json:"title"`
	CategoryName  string           `json:"categoryName"`
	UpdatedAt     string           `json:"updatedAt"`
	Comments      []domain.Comment `json:"comments"`
	Upvotes       int              `json:"upvotes"`
	CategoryID    int              `json:"categoryId"`
	Downvotes     int              `json:"downvotes"`
	Score         int              `json:"score"`
	TopicID       int              `json:"topicId"`
}

type topicPageRequest struct {
	TopicID string `url:"id"`
}

type topicPageData struct {
	User  *domain.LoggedInUser `json:"user"`
	Topic domain.Topic         `json:"topic"`
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

	requestID := &topicPageRequest{
		TopicID: topicIDStr,
	}

	backendURL, err := createURLWithParams(backendGetTopicByID, requestID)
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

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		http.Error(w, "Error with the response", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	if backendResp.StatusCode == http.StatusNotFound {
		templates.NotFoundHandler(w, r, "Topic not found", http.StatusNotFound)
		return
	}

	if backendResp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %d", backendResp.StatusCode)
		templates.NotFoundHandler(w, r, "Error loading topic", http.StatusInternalServerError)
		return
	}

	var topicData topicPageResponse
	err = helpers.DecodeBackendResponse(backendResp, &topicData)
	if err != nil {
		http.Error(w, "Error with decoding response into data struct", http.StatusInternalServerError)
		return
	}
	// log.Printf("Topic Data %v", topicData)

	// For topic related data in template
	topic := domain.Topic{
		ID:            topicData.TopicID,
		CategoryID:    topicData.CategoryID,
		Title:         topicData.Title,
		Content:       topicData.Content,
		ImagePath:     topicData.ImagePath,
		UserID:        topicData.UserID,
		CreatedAt:     topicData.CreatedAt,
		UpdatedAt:     topicData.UpdatedAt,
		UpvoteCount:   topicData.Upvotes,
		DownvoteCount: topicData.Downvotes,
		VoteScore:     topicData.Score,
		UserVote:      topicData.UserVote,
		OwnerUsername: topicData.OwnerUsername,
		CategoryName:  topicData.CategoryName,
		CategoryColor: helpers.NormalizeColor(topicData.CategoryColor),
		Comments:      topicData.Comments,
	}

	pageData := &topicPageData{
		User:  middleware.GetUserFromContext(r.Context()),
		Topic: topic,
	}

	tmpl, err := template.ParseFiles(
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
