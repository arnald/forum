// Package handlers - actions.go contains handlers for user interactive actions.
// This file handles user-initiated actions like commenting, liking, filtering posts.
// These actions require authentication and often trigger notifications.
package handlers

import (
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
	"strings"
)

// ===== COMMENT ACTIONS =====

// CreateCommentHandler processes new comment submissions (POST /create-comment).
// Requires authentication, validates input, creates comment, and triggers notifications.
func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure this is a POST request
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	// Verify user is authenticated before allowing comment creation
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extract and validate post ID from form
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	// Extract and validate comment content
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		h.BadRequest(w, r, "Comment content is required")
		return
	}

	// Create the comment in the database
	err = h.CreateComment(postID, user.ID, content)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Create notification for post author (if it's not the commenter themselves)
	// Use status-agnostic post retrieval to ensure notifications work for all posts
	post, err := h.GetPostByIDWithStatus(postID)
	if err == nil && post.UserID != user.ID {
		message := fmt.Sprintf("%s commented on your post: %s", user.Username, post.Title)
		h.CreateNotification(post.UserID, user.ID, models.NotificationTypeComment, &postID, nil, message)
	}

	// Redirect back to the post to show the new comment
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// ===== VOTING ACTIONS =====

// LikePostHandler processes post like actions (POST /like-post/{id}).
// Toggles user's like status and creates notifications for post author.
func (h *Handler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure POST method for data modification
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	// Require authentication for voting
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extract post ID from URL path (/like-post/123 -> "123")
	postIDStr := strings.TrimPrefix(r.URL.Path, "/like-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	// Toggle the like status (handles add/remove/switch logic)
	err = h.TogglePostLike(postID, user.ID, true)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Notify post author about the like (but not if they liked their own post)
	// Use status-agnostic post retrieval to ensure notifications work for all posts
	post, err := h.GetPostByIDWithStatus(postID)
	if err == nil && post.UserID != user.ID {
		message := fmt.Sprintf("%s liked your post: %s", user.Username, post.Title)
		h.CreateNotification(post.UserID, user.ID, models.NotificationTypeLike, &postID, nil, message)
	}

	// Redirect back to where the user came from (maintains browsing context)
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback to homepage
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// DislikePostHandler processes post dislike actions (POST /dislike-post/{id}).
// Similar to LikePostHandler but sets isLike=false in the voting system.
func (h *Handler) DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/dislike-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	err = h.TogglePostLike(postID, user.ID, false)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Notify post author about the dislike (but not if they disliked their own post)
	// Use status-agnostic post retrieval to ensure notifications work for all posts
	post, err := h.GetPostByIDWithStatus(postID)
	if err == nil && post.UserID != user.ID {
		message := fmt.Sprintf("%s disliked your post: %s", user.Username, post.Title)
		h.CreateNotification(post.UserID, user.ID, models.NotificationTypeDislike, &postID, nil, message)
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// LikeCommentHandler processes comment like actions (POST /like-comment/{id}).
// Similar to post voting but operates on comments instead of posts.
func (h *Handler) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentIDStr := strings.TrimPrefix(r.URL.Path, "/like-comment/")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid comment ID")
		return
	}

	err = h.ToggleCommentLike(commentID, user.ID, true)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// DislikeCommentHandler processes comment dislike actions (POST /dislike-comment/{id}).
// Identical to LikeCommentHandler but sets isLike=false.
func (h *Handler) DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentIDStr := strings.TrimPrefix(r.URL.Path, "/dislike-comment/")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid comment ID")
		return
	}

	err = h.ToggleCommentLike(commentID, user.ID, false)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// ===== FILTERING ACTIONS =====

// FilterPosts handles post filtering requests (GET /filter).
// Supports filtering by category, user's own posts, or liked posts.
// Uses query parameters to determine filter type and applies appropriate logic.
func (h *Handler) FilterPosts(w http.ResponseWriter, r *http.Request) {
	// Extract filter parameters from query string
	filter := r.URL.Query().Get("filter")         // "my-posts", "liked-posts", etc.
	categoryStr := r.URL.Query().Get("category")  // Category ID for category filtering

	var posts []models.Post
	var err error

	// Get current user - redirect to login if not authenticated
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		// User not logged in - redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse category ID if provided
	var categoryID int
	if categoryStr != "" {
		categoryID, err = strconv.Atoi(categoryStr)
		if err != nil {
			h.BadRequest(w, r, "Invalid category ID")
			return
		}
	}

	// Determine which filtering logic to apply based on parameters
	if filter == "my-posts" || filter == "liked-posts" {
		// Filter by user-specific criteria
		posts, err = h.GetPostsForUserInternalWithCategory(filter, user.ID, user, categoryID)
	} else {
		// Show posts with user-specific visibility (optionally filtered by category)
		posts, err = h.GetPostsForUserInternalWithCategory("", 0, user, categoryID)
	}

	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Get additional data needed for template rendering
	categories, _ := h.GetCategories()       // For category filter dropdown

	// Prepare data structure for template
	data := struct {
		Posts      []models.Post
		User       *models.User
		Categories []models.Category
		Filter     string  // Pass filter back to template for UI state
	}{
		Posts:      posts,
		User:       user,
		Categories: categories,
		Filter:     filter,
	}

	// Render using same template as homepage but with filtered data
	h.render(w, "home.html", data)
}

// ===== USER-SPECIFIC QUERIES =====
// Note: User-specific post queries have been moved to queries.go as GetPostsForUserInternal
// for better organization and to support role-based visibility throughout the application.
