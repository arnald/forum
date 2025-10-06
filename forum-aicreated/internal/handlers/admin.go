// Package handlers - admin.go implements administrative functionality for the forum.
// This file contains handlers for the admin panel, user role management, post moderation,
// and category management. All functions require admin-level permissions and include
// proper authorization checks to prevent unauthorized access.
package handlers

import (
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

	// Verify user has admin privileges
	// Return 404 instead of 403 to hide admin panel existence from non-admins
	if user.Role != models.RoleAdmin {
		h.NotFound(w, r)
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
		h.NotFound(w, r)
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

func (h *Handler) ApprovePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.NotFound(w, r)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/admin/approve-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	query := `UPDATE posts SET status = 'approved' WHERE id = ?`
	_, err = h.db.Exec(query, postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) RejectPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.NotFound(w, r)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/admin/reject-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid post ID")
		return
	}

	query := `UPDATE posts SET status = 'rejected' WHERE id = ?`
	_, err = h.db.Exec(query, postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

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

	if user.Role != models.RoleAdmin {
		h.NotFound(w, r)
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

	if user.Role != models.RoleAdmin {
		h.NotFound(w, r)
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
