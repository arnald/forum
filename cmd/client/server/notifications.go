package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/middleware"
)

const (
	bufferSize = 1024
)

// StreamNotifications proxies SSE stream from backend to frontend.
func (cs *ClientServer) StreamNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, cs.BackendURLs.NotificationsStreamURL(), nil)
	if err != nil {
		log.Printf("Failed to create backend request: %v", err)
		http.Error(w, "Failed to connect to notifications", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(backendReq, ip)
	for _, cookie := range r.Cookies() {
		backendReq.AddCookie(cookie)
	}

	resp, err := cs.SseClient.Do(backendReq)
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

	buf := make([]byte, bufferSize)
	for {
		select {
		case <-r.Context().Done():
			// client disconected
			return
		default:
			n, err := resp.Body.Read(buf)
			if n > 0 {
				_, writeErr := w.Write(buf[:n])
				if writeErr != nil {
					log.Printf("Write error: %v", writeErr)
					return
				}
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

// GetNotifications fetches notification list.
func (cs *ClientServer) GetNotifications(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	url := cs.BackendURLs.NotificationsListURL()
	if limit != "" {
		url = fmt.Sprintf("%s?limit=%s", url, limit)
	}

	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(backendReq, ip)

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
	err = json.NewEncoder(w).Encode(notifications)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// GetUnreadCount fetches unread notification count.
func (cs *ClientServer) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, cs.BackendURLs.UnreadCountURL(), nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(backendReq, ip)

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
	err = json.NewEncoder(w).Encode(countResp)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

// MarkNotificationAsRead marks a single notification as read.
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

	url := fmt.Sprintf("%s?id=%s", cs.BackendURLs.MarkAsReadURL(), notificationID)
	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(backendReq, ip)

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

// MarkAllNotificationsAsRead marks all notifications as read.
func (cs *ClientServer) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	backendReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, cs.BackendURLs.MarkAllAsReadURL(), nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(backendReq, ip)

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
