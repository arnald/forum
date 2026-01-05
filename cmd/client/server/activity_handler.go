package server

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

// ActivityData represents the data structure for the activity page.
type ActivityData struct {
	User             *domain.LoggedInUser
	CreatedTopics    []ActivityTopic
	LikedTopics      []ActivityTopic
	DislikedTopics   []ActivityTopic
	LikedComments    []ActivityCommentVote
	DislikedComments []ActivityCommentVote
	UserComments     []ActivityComment
}

// ActivityTopic represents a topic in the activity feed.
type ActivityTopic struct {
	ID        int
	Title     string
	CreatedAt string
}

// ActivityComment represents a comment the user made.
type ActivityComment struct {
	ID         int
	Content    string
	TopicID    int
	TopicTitle string
	CreatedAt  string
}

// ActivityCommentVote represents a comment the user liked/disliked.
type ActivityCommentVote struct {
	CommentID  int
	TopicID    int
	TopicTitle string
	CreatedAt  string
}

// ActivityPage handles requests to the activity page.
func (cs *ClientServer) ActivityPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())

	// Later, this will be replaced with a backend API call
	activityData := createStaticActivityData(user)

	tmpl, err := template.ParseFiles(
		"frontend/html/layouts/base.html",
		"frontend/html/pages/activity.html",
		"frontend/html/partials/navbar.html",
		"frontend/html/partials/footer.html",
	)
	if err != nil {
		templates.NotFoundHandler(w, r, "Failed to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", activityData)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// This will be replaced with actual backend API calls later
func createStaticActivityData(user *domain.LoggedInUser) ActivityData {
	now := time.Now()

	return ActivityData{
		User: user,
		CreatedTopics: []ActivityTopic{
			{
				ID:        1,
				Title:     "How to master Go concurrency patterns",
				CreatedAt: now.Add(-2 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
			{
				ID:        2,
				Title:     "Best practices for RESTful API design",
				CreatedAt: now.Add(-24 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
		},
		LikedTopics: []ActivityTopic{
			{
				ID:        3,
				Title:     "Understanding quantum computing basics",
				CreatedAt: now.Add(-5 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
			{
				ID:        4,
				Title:     "The future of renewable energy",
				CreatedAt: now.Add(-48 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
		},
		DislikedTopics: []ActivityTopic{
			{
				ID:        5,
				Title:     "Why cryptocurrency is overrated",
				CreatedAt: now.Add(-12 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
		},
		LikedComments: []ActivityCommentVote{
			{
				CommentID:  1,
				TopicID:    6,
				TopicTitle: "Introduction to Machine Learning",
				CreatedAt:  now.Add(-8 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
		},
		DislikedComments: []ActivityCommentVote{
			{
				CommentID:  2,
				TopicID:    7,
				TopicTitle: "The impact of social media on mental health",
				CreatedAt:  now.Add(-15 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
		},
		UserComments: []ActivityComment{
			{
				ID:         1,
				Content:    "Great explanation! This really helped me understand the concept better.",
				TopicID:    6,
				TopicTitle: "Introduction to Machine Learning",
				CreatedAt:  now.Add(-3 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
			{
				ID:         2,
				Content:    "I disagree with some points here. Have you considered the alternative approaches?",
				TopicID:    7,
				TopicTitle: "The impact of social media on mental health",
				CreatedAt:  now.Add(-36 * time.Hour).Format("Jan 02, 2006 at 15:04"),
			},
		},
	}
}
