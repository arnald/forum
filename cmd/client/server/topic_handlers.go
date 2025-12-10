package server

import "net/http"

// TopicPage handles GET requests to /topic/{id}.
func (cs *ClientServer) TopicPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
