package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/edit-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	post, err := h.GetPostByIDWithStatus(postID)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if !h.auth.HasPermission(user, "edit_own") || (!h.auth.IsOwner(user, post.UserID) && user.Role != models.RoleAdmin && user.Role != models.RoleModerator) {
		h.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		categories, err := h.GetCategories()
		if err != nil {
			h.InternalServerError(w, r, err)
			return
		}

		postCategories, err := h.GetPostCategoryIDs(postID)
		if err != nil {
			h.InternalServerError(w, r, err)
			return
		}

		data := struct {
			User           *models.User
			Post           *models.Post
			Categories     []models.Category
			PostCategories []int
		}{
			User:           user,
			Post:           post,
			Categories:     categories,
			PostCategories: postCategories,
		}

		h.render(w, "edit-post.html", data)
	case "POST":
		h.handleEditPost(w, r, user, postID)
	default:
		h.MethodNotAllowed(w, r)
	}
}

func (h *Handler) handleEditPost(w http.ResponseWriter, r *http.Request, user *models.User, postID int) {
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categories := r.Form["categories"]

	if title == "" || content == "" {
		h.BadRequest(w, r, "Title and content are required")
		return
	}

	query := `UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?`
	_, err := h.db.Exec(query, title, content, time.Now(), postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Update categories
	_, err = h.db.Exec(`DELETE FROM post_categories WHERE post_id = ?`, postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	for _, categoryStr := range categories {
		categoryID, err := strconv.Atoi(categoryStr)
		if err != nil {
			continue
		}
		_, err = h.db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
		if err != nil {
			continue
		}
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := strings.TrimPrefix(r.URL.Path, "/delete-post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	post, err := h.GetPostByIDWithStatus(postID)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if !h.auth.IsOwner(user, post.UserID) && user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.NotFound(w, r)
		return
	}

	_, err = h.db.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) EditComment(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentIDStr := strings.TrimPrefix(r.URL.Path, "/edit-comment/")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	comment, err := h.GetCommentByID(commentID)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if !h.auth.IsOwner(user, comment.UserID) && user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		data := struct {
			User    *models.User
			Comment *models.Comment
		}{
			User:    user,
			Comment: comment,
		}

		h.render(w, "edit-comment.html", data)
	case "POST":
		h.handleEditComment(w, r, user, commentID, comment.PostID)
	default:
		h.MethodNotAllowed(w, r)
	}
}

func (h *Handler) handleEditComment(w http.ResponseWriter, r *http.Request, user *models.User, commentID, postID int) {
	content := strings.TrimSpace(r.FormValue("content"))

	if content == "" {
		h.BadRequest(w, r, "Comment content is required")
		return
	}

	query := `UPDATE comments SET content = ?, updated_at = ? WHERE id = ?`
	_, err := h.db.Exec(query, content, time.Now(), commentID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentIDStr := strings.TrimPrefix(r.URL.Path, "/delete-comment/")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	comment, err := h.GetCommentByID(commentID)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if !h.auth.IsOwner(user, comment.UserID) && user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
		h.NotFound(w, r)
		return
	}

	_, err = h.db.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(comment.PostID), http.StatusSeeOther)
}

func (h *Handler) GetPostByIDWithStatus(id int) (*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.status, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`

	var post models.Post
	err := h.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath, &post.Status, &post.CreatedAt, &post.UpdatedAt, &post.Username)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (h *Handler) GetCommentByID(id int) (*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`

	var comment models.Comment
	err := h.db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.Username)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (h *Handler) GetPostCategoryIDs(postID int) ([]int, error) {
	query := `SELECT category_id FROM post_categories WHERE post_id = ?`
	rows, err := h.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryIDs []int
	for rows.Next() {
		var categoryID int
		err := rows.Scan(&categoryID)
		if err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	return categoryIDs, nil
}
