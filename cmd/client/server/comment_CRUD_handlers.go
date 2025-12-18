package server

import "net/http"

// CreateCommentPost handles POST requests to /comments/create
func (cs *ClientServer) CreateCommentPost(w http.ResponseWriter, r *http.Request) {}

// UpdateCommentPost handles POST requests to /comments/edit
func (cs *ClientServer) UpdateCommentPost(w http.ResponseWriter, r *http.Request) {}

// DeleteCommentPost handles POST requests to /comments/delete
func (cs *ClientServer) DeleteCommentPost(w http.ResponseWriter, r *http.Request) {}
