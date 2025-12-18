package server

import "net/http"

const (
	maxUploadSize = 20 << 20 // 20 MB
	uploadDir     = "frontend/static/images/uploads"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// UpdateTopicPost handles POST requests to /topics/edit.
func (cs *ClientServer) UpdateTopicPost(w http.ResponseWriter, r *http.Request) {}

// DeleteTopicPost handles POST requests to /topics/delete.
func (cs *ClientServer) DeleteTopicPost(w http.ResponseWriter, r *http.Request) {}
