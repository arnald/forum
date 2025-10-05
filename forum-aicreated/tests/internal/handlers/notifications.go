// Package handlers - notifications.go implements the real-time notification system.
// This file handles creating, retrieving, and managing user notifications for
// forum interactions like likes, comments, and other user activities.
// Notifications help users stay engaged by alerting them to relevant activity.
package handlers

import (
	"forum/internal/models"
	"net/http"
	"strconv"
)

// ===== NOTIFICATION CREATION =====

// CreateNotification creates a new notification for a user.
// This is called whenever an action occurs that should notify another user.
// Prevents self-notifications (users don't get notified about their own actions).
func (h *Handler) CreateNotification(userID, actorID int, notifType models.NotificationType, postID, commentID *int, message string) error {
	// Don't create notifications for users' own actions
	// This prevents spam notifications when users interact with their own content
	if userID == actorID {
		return nil
	}

	// Insert notification into database
	// postID and commentID are optional (can be nil) depending on notification type
	query := `INSERT INTO notifications (user_id, actor_id, type, post_id, comment_id, message) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := h.db.Exec(query, userID, actorID, notifType, postID, commentID, message)
	return err
}

// ===== NOTIFICATION RETRIEVAL =====

// GetUserNotifications retrieves all notifications for a specific user.
// Returns notifications with actor usernames for display purposes.
// Limited to 50 most recent notifications, ordered newest first.
func (h *Handler) GetUserNotifications(userID int) ([]models.Notification, error) {
	// JOIN with users table to get actor username for display
	// Orders by creation time (newest first) and limits to prevent excessive data
	query := `
		SELECT n.id, n.user_id, n.actor_id, n.type, n.post_id, n.comment_id, n.message, n.read, n.created_at, u.username
		FROM notifications n
		JOIN users u ON n.actor_id = u.id                  -- Get actor username
		WHERE n.user_id = ?                                -- Filter by recipient
		ORDER BY n.created_at DESC                         -- Newest first
		LIMIT 50                                           -- Prevent excessive data
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		err := rows.Scan(&notif.ID, &notif.UserID, &notif.ActorID, &notif.Type, &notif.PostID, &notif.CommentID, &notif.Message, &notif.Read, &notif.CreatedAt, &notif.ActorName)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// ===== NOTIFICATION MANAGEMENT =====

// MarkNotificationRead marks a single notification as read.
// Used when user clicks on or views a specific notification.
func (h *Handler) MarkNotificationRead(notificationID int) error {
	query := `UPDATE notifications SET read = TRUE WHERE id = ?`
	_, err := h.db.Exec(query, notificationID)
	return err
}

// GetUnreadNotificationCount returns the number of unread notifications for a user.
// Used to display notification badges/counters in the UI.
func (h *Handler) GetUnreadNotificationCount(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read = FALSE`
	err := h.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// ===== HTTP HANDLERS =====

// Notifications displays the user's notification center (GET /notifications).
// Shows all notifications with read/unread status and allows interaction.
func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	// Require authentication to view notifications
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve all notifications for this user
	notifications, err := h.GetUserNotifications(user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Prepare data for template rendering
	data := struct {
		User          *models.User
		Notifications []models.Notification
	}{
		User:          user,
		Notifications: notifications,
	}

	// Render notifications page
	h.render(w, "notifications.html", data)
}

// MarkNotificationReadHandler marks a single notification as read (POST /mark-read).
// Used when user clicks on individual notifications to mark them as seen.
func (h *Handler) MarkNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure POST method for data modification
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	// Require authentication
	_, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extract notification ID from form
	notificationIDStr := r.FormValue("notification_id")
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		h.BadRequest(w, r, "Invalid notification ID")
		return
	}

	// Mark the notification as read
	err = h.MarkNotificationRead(notificationID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Redirect back to where user came from (usually notifications page)
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/notifications" // Fallback to notifications page
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// MarkAllNotificationsRead marks all of a user's notifications as read (POST /mark-all-read).
// Convenient bulk operation for users with many unread notifications.
func (h *Handler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	// Ensure POST method for data modification
	if r.Method != "POST" {
		h.MethodNotAllowed(w, r)
		return
	}

	// Require authentication
	user, err := h.auth.GetUserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Mark all notifications for this user as read
	query := `UPDATE notifications SET read = TRUE WHERE user_id = ?`
	_, err = h.db.Exec(query, user.ID)
	if err != nil {
		h.InternalServerError(w, r, err)
		return
	}

	// Redirect back to notifications page to show updated status
	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}
