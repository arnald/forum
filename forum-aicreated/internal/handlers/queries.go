// Package handlers - queries.go contains all database query methods.
// This file implements the data access layer, providing methods to retrieve
// and manipulate forum data. It uses prepared statements for security and
// efficient JOIN operations for performance.
package handlers

import (
	"database/sql"
	"forum/internal/models"
)

// ===== POST QUERIES =====

// GetPosts retrieves posts from the database with role-based visibility.
// Returns posts with computed like/dislike counts and author usernames.
// Uses JOIN operations to minimize database round trips.
// - Anonymous users: only approved posts
// - Regular users: approved posts + their own posts (any status)
// - Admins/Moderators: all posts regardless of status
func (h *Handler) GetPosts(filter string) ([]models.Post, error) {
	return h.GetPostsForUserInternal("", 0, nil)
}

// GetPostsWithUser retrieves posts with user-specific visibility rules
func (h *Handler) GetPostsWithUser(user *models.User) ([]models.Post, error) {
	return h.GetPostsForUserInternal("", 0, user)
}

// GetPostsForUserInternal retrieves posts based on filter and user permissions
// Supports categoryID for filtering by specific category
func (h *Handler) GetPostsForUserInternal(filter string, userID int, currentUser *models.User) ([]models.Post, error) {
	return h.GetPostsForUserInternalWithCategory(filter, userID, currentUser, 0)
}

// GetPostsForUserInternalWithCategory retrieves posts based on filter, user permissions, and optional category
func (h *Handler) GetPostsForUserInternalWithCategory(filter string, userID int, currentUser *models.User, categoryID int) ([]models.Post, error) {
	// Complex query that JOINs posts, users, and aggregates likes/dislikes
	// COALESCE ensures we get 0 instead of NULL for posts with no votes
	baseQuery := `
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.username,
		       COALESCE(SUM(CASE WHEN pl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes,
		       p.status
		FROM posts p
		JOIN users u ON p.user_id = u.id                    -- Get author username
		LEFT JOIN post_likes pl ON p.id = pl.post_id        -- Get vote counts (LEFT JOIN for posts with no votes)`

	// Add category JOIN if filtering by category
	if categoryID > 0 {
		baseQuery += " JOIN post_categories pc ON p.id = pc.post_id"
	}

	baseQuery += " WHERE "

	// Build WHERE clause based on user role and permissions
	var whereClause string
	var args []interface{}

	if currentUser == nil {
		// Anonymous users: only approved posts
		whereClause = "p.status = 'approved'"
	} else if currentUser.Role == models.RoleAdmin || currentUser.Role == models.RoleModerator {
		// Admins and moderators: see all posts
		whereClause = "1=1" // No status restriction
	} else {
		// Regular users: approved posts + their own posts
		whereClause = "(p.status = 'approved' OR p.user_id = ?)"
		args = append(args, currentUser.ID)
	}

	// Add category filter
	if categoryID > 0 {
		whereClause += " AND pc.category_id = ?"
		args = append(args, categoryID)
	}

	// Add user-specific filters
	if filter == "my-posts" && currentUser != nil {
		whereClause += " AND p.user_id = ?"
		args = append(args, currentUser.ID)
	} else if filter == "liked-posts" && currentUser != nil {
		baseQuery += " JOIN post_likes pl2 ON p.id = pl2.post_id AND pl2.user_id = ? AND pl2.is_like = 1"
		args = append(args, currentUser.ID)
	}

	finalQuery := baseQuery + whereClause + " GROUP BY p.id ORDER BY p.created_at DESC"

	// Execute the query
	rows, err := h.db.Query(finalQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process each row into a Post struct
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath, &post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes, &post.Status)
		if err != nil {
			return nil, err
		}

		// Fetch categories for this post (separate query for many-to-many relationship)
		categories, err := h.GetPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostByID retrieves a single post by its ID for anonymous users.
// Only returns approved posts to prevent access to pending/rejected content.
func (h *Handler) GetPostByID(id int) (*models.Post, error) {
	return h.GetPostByIDWithUser(id, nil)
}

// GetPostByIDWithUser retrieves a single post by its ID with user-specific visibility.
// Applies the same role-based visibility rules as GetPostsForUserInternal:
// - Anonymous users: only approved posts
// - Regular users: approved posts + their own posts (any status)
// - Admins/Moderators: all posts regardless of status
func (h *Handler) GetPostByIDWithUser(id int, currentUser *models.User) (*models.Post, error) {
	// Same complex query as GetPosts but filtered by specific ID
	baseQuery := `
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.username,
		       COALESCE(SUM(CASE WHEN pl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes,
		       p.status
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		WHERE p.id = ? AND `

	// Build WHERE clause based on user role and permissions
	var whereClause string
	var args []interface{}
	args = append(args, id) // First argument is always the post ID

	if currentUser == nil {
		// Anonymous users: only approved posts
		whereClause = "p.status = 'approved'"
	} else if currentUser.Role == models.RoleAdmin || currentUser.Role == models.RoleModerator {
		// Admins and moderators: see all posts
		whereClause = "1=1" // No status restriction
	} else {
		// Regular users: approved posts + their own posts
		whereClause = "(p.status = 'approved' OR p.user_id = ?)"
		args = append(args, currentUser.ID)
	}

	finalQuery := baseQuery + whereClause + " GROUP BY p.id"

	var post models.Post
	// QueryRow for single result - will return sql.ErrNoRows if not found
	err := h.db.QueryRow(finalQuery, args...).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath, &post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes, &post.Status)
	if err != nil {
		return nil, err
	}

	// Fetch categories for this post
	categories, err := h.GetPostCategories(post.ID)
	if err != nil {
		return nil, err
	}
	post.Categories = categories

	return &post, nil
}

// ===== CATEGORY QUERIES =====

// GetPostCategories retrieves all category names for a specific post.
// Handles the many-to-many relationship between posts and categories.
func (h *Handler) GetPostCategories(postID int) ([]string, error) {
	// JOIN categories with post_categories junction table
	query := `
		SELECT c.name
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id    -- Join through junction table
		WHERE pc.post_id = ?                                -- Filter by specific post
	`

	rows, err := h.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect category names into slice
	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategories retrieves all available categories for post organization.
// Used in forms where users can select categories for their posts.
func (h *Handler) GetCategories() ([]models.Category, error) {
	// Simple query - categories table has no relationships to worry about
	query := `SELECT id, name FROM categories ORDER BY name`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// ===== COMMENT QUERIES =====

// GetCommentsByPostID retrieves all comments for a specific post.
// Similar to post queries, includes computed like/dislike counts and usernames.
// Comments are ordered chronologically (oldest first) for natural reading flow.
func (h *Handler) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	// Complex query similar to posts but for comments
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username,
		       COALESCE(SUM(CASE WHEN cl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN cl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id                    -- Get commenter username
		LEFT JOIN comment_likes cl ON c.id = cl.comment_id  -- Get vote counts
		WHERE c.post_id = ?                                 -- Filter by specific post
		GROUP BY c.id
		ORDER BY c.created_at ASC                           -- Chronological order (oldest first)
	`

	rows, err := h.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// ===== DATA MODIFICATION QUERIES =====

// CreateComment adds a new comment to a post.
// Simple INSERT operation with foreign key relationships.
func (h *Handler) CreateComment(postID, userID int, content string) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	_, err := h.db.Exec(query, postID, userID, content)
	return err
}

// ===== VOTING SYSTEM QUERIES =====

// TogglePostLike handles like/dislike functionality for posts.
// Implements complex logic: INSERT if no vote, DELETE if same vote, UPDATE if different vote.
// This allows users to like, dislike, or remove their vote entirely.
func (h *Handler) TogglePostLike(postID, userID int, isLike bool) error {
	// First, check if user has already voted on this post
	var existingLike sql.NullBool
	query := `SELECT is_like FROM post_likes WHERE post_id = ? AND user_id = ?`
	err := h.db.QueryRow(query, postID, userID).Scan(&existingLike)

	if err == sql.ErrNoRows {
		// No existing vote - create new one
		query = `INSERT INTO post_likes (post_id, user_id, is_like) VALUES (?, ?, ?)`
		_, err = h.db.Exec(query, postID, userID, isLike)
		return err
	} else if err != nil {
		return err
	}

	// Check if user is clicking the same vote type (like when already liked)
	if existingLike.Valid && existingLike.Bool == isLike {
		// Remove the vote (toggle off)
		query = `DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`
		_, err = h.db.Exec(query, postID, userID)
		return err
	}

	// User is switching from like to dislike or vice versa
	query = `UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?`
	_, err = h.db.Exec(query, isLike, postID, userID)
	return err
}

// ToggleCommentLike handles like/dislike functionality for comments.
// Identical logic to TogglePostLike but operates on comment_likes table.
func (h *Handler) ToggleCommentLike(commentID, userID int, isLike bool) error {
	var existingLike sql.NullBool
	query := `SELECT is_like FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	err := h.db.QueryRow(query, commentID, userID).Scan(&existingLike)

	if err == sql.ErrNoRows {
		query = `INSERT INTO comment_likes (comment_id, user_id, is_like) VALUES (?, ?, ?)`
		_, err = h.db.Exec(query, commentID, userID, isLike)
		return err
	} else if err != nil {
		return err
	}

	if existingLike.Valid && existingLike.Bool == isLike {
		query = `DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?`
		_, err = h.db.Exec(query, commentID, userID)
		return err
	}

	query = `UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?`
	_, err = h.db.Exec(query, isLike, commentID, userID)
	return err
}

// ===== FILTERING QUERIES =====

// GetPostsByCategory retrieves posts filtered by a specific category.
// Similar to GetPosts but adds category filtering through the junction table.
func (h *Handler) GetPostsByCategory(categoryID int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.username,
		       COALESCE(SUM(CASE WHEN pl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		WHERE pc.category_id = ? AND p.status = 'approved'
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := h.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath, &post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}

		categories, err := h.GetPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}
