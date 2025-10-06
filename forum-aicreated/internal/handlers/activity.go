// Package handlers - activity.go implements user activity tracking and display.
// This file provides functionality to retrieve and display user's forum activity
// including posts they've created, posts they've liked, and comments they've made.
// Used for user profiles and personal activity dashboards.
package handlers

import (
	"forum/internal/models"
	"net/http"
)

// ===== ACTIVITY DATA STRUCTURES =====

// UserActivity aggregates all user activity data for profile/dashboard display
type UserActivity struct {
	CreatedPosts []models.Post     `json:"created_posts"` // Posts authored by the user
	LikedPosts   []models.Post     `json:"liked_posts"`   // Posts the user has liked
	Comments     []ActivityComment `json:"comments"`      // Comments made by the user
}

// ActivityComment extends the base Comment model with post title for context
// This allows displaying "User commented on: Post Title" in activity feeds
type ActivityComment struct {
	models.Comment
	PostTitle string `json:"post_title"` // Title of the post that was commented on
}

// ===== ACTIVITY RETRIEVAL FUNCTIONS =====

// GetUserCreatedPosts retrieves all posts created by a specific user
// Uses the existing post filtering system with user context for proper visibility
func (h *Handler) GetUserCreatedPosts(userID int) ([]models.Post, error) {
	// Create a mock user object for the filtering system
	// This ensures proper visibility rules are applied
	user := &models.User{ID: userID}
	return h.GetPostsForUserInternal("my-posts", userID, user)
}

// GetUserLikedPosts retrieves all posts that a user has liked
// Uses the existing post filtering system with proper visibility rules
func (h *Handler) GetUserLikedPosts(userID int) ([]models.Post, error) {
	// Create a mock user object for the filtering system
	user := &models.User{ID: userID}
	return h.GetPostsForUserInternal("liked-posts", userID, user)
}

// GetUserComments retrieves all comments made by a specific user with post context
// Includes post titles and like/dislike counts for comprehensive activity display
// Limited to 50 most recent comments to prevent excessive data loading
func (h *Handler) GetUserComments(userID int) ([]ActivityComment, error) {
	// Complex query that joins comments with users, posts, and aggregates likes
	// Includes post title so we can show "User commented on: Post Title"
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.username, p.title as post_title,
		       COALESCE(SUM(CASE WHEN cl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN cl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id                   -- Get commenter username
		JOIN posts p ON c.post_id = p.id                   -- Get post title for context
		LEFT JOIN comment_likes cl ON c.id = cl.comment_id -- Get like/dislike counts
		WHERE c.user_id = ?                                -- Filter by specific user
		GROUP BY c.id                                      -- Group to aggregate likes
		ORDER BY c.created_at DESC                         -- Newest first
		LIMIT 50                                           -- Prevent excessive data
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []ActivityComment
	for rows.Next() {
		var comment ActivityComment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt, &comment.Username, &comment.PostTitle,
			&comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// ===== ACTIVITY PAGE HANDLER =====

// Activity handles the user activity page display (GET /activity)
// Shows a comprehensive view of user's forum activity including posts, likes, and comments
// Requires authentication - redirects to login if user not logged in
func (h *Handler) Activity(w http.ResponseWriter, r *http.Request) {
	// Verify user is authenticated
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve all posts created by the user
	createdPosts, err := h.GetUserCreatedPosts(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Retrieve all posts the user has liked
	likedPosts, err := h.GetUserLikedPosts(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Retrieve all comments made by the user
	comments, err := h.GetUserComments(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Aggregate all activity data for template rendering
	activity := UserActivity{
		CreatedPosts: createdPosts,
		LikedPosts:   likedPosts,
		Comments:     comments,
	}

	// Prepare data structure for template rendering
	// Includes both user context and activity data
	data := struct {
		User     *models.User
		Activity UserActivity
	}{
		User:     user,
		Activity: activity,
	}

	// Render the activity page template
	h.render(w, "activity.html", data)
}
