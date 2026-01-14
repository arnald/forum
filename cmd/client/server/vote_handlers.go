package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/arnald/forum/cmd/client/helpers"
)

type voteCountsResponse struct {
	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
	Score     int `json:"score"`
}

func (cs *ClientServer) proxyVoteRequest(w http.ResponseWriter, r *http.Request, backendURL string, method string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Request proxyVoteRequest body: %s", string(body))

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, method, backendURL, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	for _, cookie := range r.Cookies() {
		log.Printf("Forwarding cookie: %s=%s", cookie.Name, cookie.Value)
		httpReq.AddCookie(cookie)
	}

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("Error making request to backend: %v", err)
		http.Error(w, "Error communicating with backend", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	respBody, err := io.ReadAll(backendResp.Body)
	if err != nil {
		http.Error(w, "Error reading backend response", http.StatusInternalServerError)
		return
	}
	log.Printf("Response: %s", string(respBody))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(backendResp.StatusCode)
	_, err = w.Write(respBody)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// CastVote proxies the vote casting request to the backend.
func (cs *ClientServer) CastVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cs.proxyVoteRequest(w, r, backendCastVote, http.MethodPost)
}

// DeleteVote proxies the vote deletion request to the backend.
func (cs *ClientServer) DeleteVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cs.proxyVoteRequest(w, r, backendDeleteVote, http.MethodDelete)
}

// GetVoteCounts gets the current vote counts for a topic or comment.
func (cs *ClientServer) GetVoteCounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topicIDStr := r.URL.Query().Get("topic_id")
	commentIDStr := r.URL.Query().Get("comment_id")

	if topicIDStr == "" && commentIDStr == "" {
		http.Error(w, "Either topic_id or comment_id must be provided", http.StatusBadRequest)
		return
	}

	backendURL := backendGetVoteCounts + "?"
	if topicIDStr != "" {
		backendURL += "topic_id=" + topicIDStr
	} else {
		backendURL += "comment_id=" + commentIDStr
	}
	log.Printf("Backend GetVoteCounts URL %v", backendURL)

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendURL, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		httpReq.AddCookie(cookie)
	}

	backendResp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("Error making request to backend: %v", err)
		http.Error(w, "Error communicating with backend", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	if backendResp.StatusCode != http.StatusOK {
		log.Printf("Backend returned status: %d", backendResp.StatusCode)
		http.Error(w, "Error getting vote counts", backendResp.StatusCode)
		return
	}

	var counts voteCountsResponse
	err = helpers.DecodeBackendResponse(backendResp, &counts)
	if err != nil {
		http.Error(w, "Error decoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Data voteCountsResponse `json:"data"`
	}{
		Data: counts,
	})
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
