// Package handlers - admin.go implements administrative functionality for the forum.
// This file contains handlers for the admin panel, user role management, post moderation,
// and category management. All functions require admin-level permissions and include
// proper authorization checks to prevent unauthorized access.
package handlers

import (
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
	"strings"
)

// ===== ADMIN PANEL HANDLERS =====

// AdminPanel displays the main administrative dashboard (GET /admin)
// Shows user management, post moderation, category management, and reports
// Requires admin role - returns 404 for non-admin users for security
func (h *Handler) AdminPanel(w http.ResponseWriter, r *http.Request) {
	// Verify user is authenticated
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verify user has admin or moderator privileges
	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.Forbidden(w, r, "You don't have permission to access the admin panel. Only admins and moderators can access this area.")
		return
	}

	users, err := h.GetAllUsers()
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	categories, err := h.GetCategories()
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	pendingPosts, err := h.GetPendingPosts()
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	reports, err := h.GetPendingReports()
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	data := struct {
		User         *models.User
		Users        []models.User
		Categories   []models.Category
		PendingPosts []models.Post
		Reports      []models.Report
	}{
		User:         user,
		Users:        users,
		Categories:   categories,
		PendingPosts: pendingPosts,
		Reports:      reports,
	}

	h.render(w, "admin.html", data)
}

func (h *Handler) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, email, username, role, provider, created_at FROM users ORDER BY created_at DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.Provider, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (h *Handler) GetPendingPosts() ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.status = 'pending'
		ORDER BY p.created_at DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath, &post.CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (h *Handler) GetPendingReports() ([]models.Report, error) {
	query := `
		SELECT r.id, r.reporter_id, r.post_id, r.comment_id, r.reason, r.description, r.status, r.created_at, u.username
		FROM reports r
		JOIN users u ON r.reporter_id = u.id
		WHERE r.status = 'pending'
		ORDER BY r.created_at DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		err := rows.Scan(&report.ID, &report.ReporterID, &report.PostID, &report.CommentID, &report.Reason, &report.Description, &report.Status, &report.CreatedAt, &report.ReporterName)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != models.RoleAdmin {
		h.Forbidden(w, r, "You don't have permission to change user roles. Only admins can modify user roles.")
		return
	}

	userIDStr := r.FormValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid user ID")
		return
	}

	newRole := models.UserRole(r.FormValue("role"))
	if newRole != models.RoleUser && newRole != models.RoleModerator && newRole != models.RoleAdmin {
		h.BadRequest(w, r, "Invalid role")
		return
	}

	err = h.auth.UpdateUserRole(userID, newRole)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// ===== POST MODERATION HANDLERS =====

// ApprovePost handles post approval by admins and moderators (POST /admin/approve-post/{id})
// Changes post status from 'pending' to 'approved', making it visible to all users
// Sends notification to post author about the approval decision
// Enhanced to include user notifications for better user experience
func (h *Handler) ApprovePost(w http.ResponseWriter, r *http.Request) {
	// Ensure POST method for data modification
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	// Verify user is authenticated
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verify user has moderation privileges (admin OR moderator)
	// This allows both roles to approve posts for efficient content moderation
	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.Forbidden(w, r, "You don't have permission to approve posts. Only admins and moderators can approve posts.")
		return
	}

	// Extract post ID from URL path
	postIDStr := strings.TrimPrefix(r.URL.Path, "/admin/approve-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	// Get post details before approval for notification purposes
	// Uses GetPostByIDWithStatus to retrieve posts regardless of current status
	post, err := h.GetPostByIDWithStatus(postID)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	// Update post status to approved in database
	// This makes the post visible to all users in the main forum feed
	query := `UPDATE posts SET status = 'approved' WHERE id = ?`
	_, err = h.db.Exec(query, postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Notify post author about approval decision
	// Provides feedback to users about their content status for better UX
	message := fmt.Sprintf("Your post '%s' has been approved by %s", post.Title, user.Username)
	h.CreateNotification(post.UserID, user.ID, models.NotificationTypePostApproved, &postID, nil, message)

	// Redirect back to admin panel to continue moderation workflow
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// RejectPost handles post rejection by admins and moderators (POST /admin/reject-post/{id})
// Changes post status from 'pending' to 'rejected', keeping it hidden from public view
// Post remains visible to author in their "My Posts" section with rejected status
// Sends notification to post author about the rejection decision with moderator details
func (h *Handler) RejectPost(w http.ResponseWriter, r *http.Request) {
	// Ensure POST method for data modification
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	// Verify user is authenticated
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verify user has moderation privileges (admin OR moderator)
	// Both roles can reject posts to maintain content quality standards
	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.Forbidden(w, r, "You don't have permission to reject posts. Only admins and moderators can reject posts.")
		return
	}

	// Extract post ID from URL path
	postIDStr := strings.TrimPrefix(r.URL.Path, "/admin/reject-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	// Get post details before rejection for notification purposes
	// Uses GetPostByIDWithStatus to handle posts in any status (pending, etc.)
	post, err := h.GetPostByIDWithStatus(postID)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	// Update post status to rejected in database
	// Rejected posts remain hidden from public but visible to author
	query := `UPDATE posts SET status = 'rejected' WHERE id = ?`
	_, err = h.db.Exec(query, postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Notify post author about rejection decision
	// Important for transparency - users should know why their content isn't visible
	message := fmt.Sprintf("Your post '%s' has been rejected by %s", post.Title, user.Username)
	h.CreateNotification(post.UserID, user.ID, models.NotificationTypePostRejected, &postID, nil, message)

	// Redirect back to admin panel to continue moderation workflow
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Allow both admins and moderators to create categories for content organization
	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.Forbidden(w, r, "You don't have permission to create categories. Only admins and moderators can create categories.")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		h.BadRequest(w, r, "Category name is required")
		return
	}

	query := `INSERT INTO categories (name) VALUES (?)`
	_, err = h.db.Exec(query, name)
	if err != nil {
		h.BadRequest(w, r, "Category already exists or invalid name")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Allow both admins and moderators to delete categories for content organization
	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.Forbidden(w, r, "You don't have permission to delete categories. Only admins and moderators can delete categories.")
		return
	}

	categoryIDStr := strings.TrimPrefix(r.URL.Path, "/admin/delete-category/")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid category ID")
		return
	}

	query := `DELETE FROM categories WHERE id = ?`
	_, err = h.db.Exec(query, categoryID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// ===== BULK MODERATION ACTIONS =====

// ApproveModeratorPosts automatically approves all pending posts created by users who are now moderators or admins
// This is useful when users are promoted to moderator/admin roles and their old pending posts should be approved
func (h *Handler) ApproveModeratorPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only admins can perform bulk operations
	if user.Role != models.RoleAdmin {
		h.Forbidden(w, r, "You don't have permission to perform bulk operations. Only admins can do this.")
		return
	}

	// Update all pending posts created by moderators or admins to approved status
	query := `
		UPDATE posts
		SET status = 'approved'
		WHERE status = 'pending'
		AND user_id IN (SELECT id FROM users WHERE role IN ('moderator', 'admin'))
	`
	result, err := h.db.Exec(query)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Get number of rows affected for user feedback
	rowsAffected, _ := result.RowsAffected()

	// Log or notify about the bulk action
	if rowsAffected > 0 {
		// You could add a success message here via session flash or query parameter
		// For now, just redirect back to admin panel
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// BulkApproveAllPending approves ALL pending posts (use with caution)
// This is a nuclear option for clearing the pending queue
func (h *Handler) BulkApproveAllPending(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only admins can perform bulk operations
	if user.Role != models.RoleAdmin {
		h.Forbidden(w, r, "You don't have permission to perform bulk operations. Only admins can do this.")
		return
	}

	// Approve all pending posts
	query := `UPDATE posts SET status = 'approved' WHERE status = 'pending'`
	_, err = h.db.Exec(query)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// BulkRejectAllPending rejects ALL pending posts (use with extreme caution)
func (h *Handler) BulkRejectAllPending(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Only admins can perform bulk operations
	if user.Role != models.RoleAdmin {
		h.Forbidden(w, r, "You don't have permission to perform bulk operations. Only admins can do this.")
		return
	}

	// Reject all pending posts
	query := `UPDATE posts SET status = 'rejected' WHERE status = 'pending'`
	_, err = h.db.Exec(query)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
