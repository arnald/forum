package handlers

import (
	"forum/internal/models"
	"net/http"
)

type UserActivity struct {
	CreatedPosts []models.Post     `json:"created_posts"`
	LikedPosts   []models.Post     `json:"liked_posts"`
	Comments     []ActivityComment `json:"comments"`
}

type ActivityComment struct {
	models.Comment
	PostTitle string `json:"post_title"`
}

func (h *Handler) GetUserCreatedPosts(userID int) ([]models.Post, error) {
	return h.GetPostsForUser("my-posts", userID)
}

func (h *Handler) GetUserLikedPosts(userID int) ([]models.Post, error) {
	return h.GetPostsForUser("liked-posts", userID)
}

func (h *Handler) GetUserComments(userID int) ([]ActivityComment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.username, p.title as post_title,
		       COALESCE(SUM(CASE WHEN cl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN cl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		JOIN posts p ON c.post_id = p.id
		LEFT JOIN comment_likes cl ON c.id = cl.comment_id
		WHERE c.user_id = ?
		GROUP BY c.id
		ORDER BY c.created_at DESC
		LIMIT 50
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

func (h *Handler) Activity(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	createdPosts, err := h.GetUserCreatedPosts(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	likedPosts, err := h.GetUserLikedPosts(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	comments, err := h.GetUserComments(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	activity := UserActivity{
		CreatedPosts: createdPosts,
		LikedPosts:   likedPosts,
		Comments:     comments,
	}

	data := struct {
		User     *models.User
		Activity UserActivity
	}{
		User:     user,
		Activity: activity,
	}

	h.render(w, "activity.html", data)
}
