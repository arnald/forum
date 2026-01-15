package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arnald/forum/cmd/client/domain"
)

const (
	backendNotificationsStream = "http://localhost:8080/api/v1/notifications/stream"
	backendNotificationsList   = "http://localhost:8080/api/v1/notifications"
	backendUnreadCount         = "http://localhost:8080/api/v1/notifications/unread-count"
	backendMarkAsRead          = "http://localhost:8080/api/v1/notifications/mark-read"
	backendMarkAllAsRead       = "http://localhost:8080/api/v1/notifications/mark-all-read"
)

// StreamNotifications proxies SSE stream from backend to frontend
func (cs *ClientServer) StreamNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, backendNotificationsStream, nil)
	if err != nil {
		log.Printf("Failed to create backend request: %v", err)
		http.Error(w, "Failed to connect to notifications", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		backendReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(backendReq)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		http.Error(w, "Failed to connect to notifications", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %d", resp.StatusCode)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Stream data from backend to client
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	buf := make([]byte, 1024)
	for {
		select {
		case <-r.Context().Done():
			// Client disconnected (browser refresh, tab close, etc.)
			log.Println("Client disconnected from SSE stream")
			return
		default:
			n, err := resp.Body.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
				flusher.Flush()
			}
			if err != nil {
				// Backend closed connection or error occurred
				if err != io.EOF {
					log.Printf("Stream read error: %v", err)
				}
				return
			}
		}
	}
}

// GetNotifications fetches notification list
func (cs *ClientServer) GetNotifications(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	url := backendNotificationsList
	if limit != "" {
		url = fmt.Sprintf("%s?limit=%s", url, limit)
	}

	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		backendReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(backendReq)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch notifications", resp.StatusCode)
		return
	}

	var notifications []*domain.Notification
	err = json.NewDecoder(resp.Body).Decode(&notifications)
	if err != nil {
		log.Printf("Failed to decode response: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// GetUnreadCount fetches unread notification count
func (cs *ClientServer) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, backendUnreadCount, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		backendReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(backendReq)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		http.Error(w, "Failed to fetch count", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch count", resp.StatusCode)
		return
	}

	var countResp domain.UnreadCountResponse
	err = json.NewDecoder(resp.Body).Decode(&countResp)
	if err != nil {
		log.Printf("Failed to decode response: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countResp)
}

// MarkNotificationAsRead marks a single notification as read
func (cs *ClientServer) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	notificationID := r.URL.Query().Get("id")
	if notificationID == "" {
		http.Error(w, "Missing notification ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s?id=%s", backendMarkAsRead, notificationID)
	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		backendReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(backendReq)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		http.Error(w, "Failed to mark as read", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to mark as read", resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// MarkAllNotificationsAsRead marks all notifications as read
func (cs *ClientServer) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, backendMarkAllAsRead, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		backendReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(backendReq)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		http.Error(w, "Failed to mark all as read", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to mark all as read", resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
