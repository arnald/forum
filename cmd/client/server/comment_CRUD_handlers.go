package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/arnald/forum/cmd/client/helpers/templates"
)

type createCommentRequest struct {
	Content string `json:"content"`
	TopicID int    `json:"topicId"`
}

type updateCommentRequest struct {
	Content string `json:"content"`
	ID      int    `json:"id"`
}

// CreateCommentPost handles POST requests to /comments/create.
func (cs *ClientServer) CreateCommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	topicIDStr := r.FormValue("topic_id")
	content := r.FormValue("content")

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		log.Printf("Invalid topic ID: %v", err)
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	if content == "" {
		log.Printf("Empty comment content")
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	createRequest := &createCommentRequest{
		TopicID: topicID,
		Content: content,
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := cs.newRequestWithCookies(ctx, http.MethodPost, backendCreateComment, createRequest, r)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		templates.NotFoundHandler(w, r, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Backend returned error: %s", string(body))
		templates.NotFoundHandler(w, r, "Failed to create comment", resp.StatusCode)
		return
	}

	http.Redirect(w, r, "/topic/"+topicIDStr, http.StatusSeeOther)
}

// UpdateCommentPost handles POST requests to /comments/edit.
func (cs *ClientServer) UpdateCommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	content := r.FormValue("content")
	topicIDStr := r.FormValue("topic_id")

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Printf("Invalid comment ID: %v", err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if content == "" {
		log.Printf("Empty comment content")
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	updateRequest := &updateCommentRequest{
		ID:      commentID,
		Content: content,
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := cs.newRequestWithCookies(ctx, http.MethodPut, backendUpdateComment, updateRequest, r)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		templates.NotFoundHandler(w, r, "Failed to update comment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Backend returned error: %s", string(body))
		templates.NotFoundHandler(w, r, "Failed to update comment", resp.StatusCode)
		return
	}

	if topicIDStr == "" {
		log.Printf("No topic_id provided, redirecting to topics list")
		http.Redirect(w, r, "/topics", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/topic/"+topicIDStr, http.StatusSeeOther)
}

// DeleteCommentPost handles POST requests to /comments/delete.
func (cs *ClientServer) DeleteCommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	topicIDStr := r.FormValue("topic_id")

	_, err = strconv.Atoi(commentIDStr)
	if err != nil {
		log.Printf("Invalid comment ID: %v", err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	deleteURL := backendDeleteComment + "?id=" + commentIDStr

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		httpReq.AddCookie(cookie)
	}

	resp, err := cs.HTTPClient.Do(httpReq)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		templates.NotFoundHandler(w, r, "Failed to delete comment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Backend returned error: %s", string(body))
		templates.NotFoundHandler(w, r, "Failed to delete comment", resp.StatusCode)
		return
	}

	if topicIDStr == "" {
		log.Printf("No topic_id provided, redirecting to topics list")
		http.Redirect(w, r, "/topics", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/topic/"+topicIDStr, http.StatusSeeOther)
}
