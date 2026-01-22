package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/arnald/forum/cmd/client/domain"
	"github.com/arnald/forum/cmd/client/helpers"
	"github.com/arnald/forum/cmd/client/helpers/templates"
	"github.com/arnald/forum/cmd/client/middleware"
)

const (
	maxUploadSize = 20 << 20 // 20 MB
	uploadDir     = "frontend/static/images/uploads"
	uploadDirPerm = 0o750
)

type createTopicRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ImagePath   string `json:"imagePath"`
	CategoryIDs []int  `json:"categoryIds"`
}

type updateTopicRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ImagePath   string `json:"imagePath"`
	CategoryIDs []int  `json:"categoryIds"`
	TopicID     int    `json:"topicId"`
}

type createPostData struct {
	Categories []domain.Category
}

// CreateTopicPage handles GET requests to /topics/create - shows the form.
func (cs *ClientServer) CreateTopicPage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	categoriesHTTPReq, err := http.NewRequestWithContext(ctx, http.MethodGet, backendGetCategoriesDomain(), nil)
	if err != nil {
		log.Printf("Error creating categories request: %v", err)
		templates.NotFoundHandler(w, r, "Error creating categories request", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(categoriesHTTPReq, ip)

	for _, cookie := range r.Cookies() {
		categoriesHTTPReq.AddCookie(cookie)
	}

	categoriesResp, err := cs.HTTPClient.Do(categoriesHTTPReq)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		templates.NotFoundHandler(w, r, "Error fetching categories", http.StatusInternalServerError)
		return
	}
	defer categoriesResp.Body.Close()

	if categoriesResp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch categories, status: %d", categoriesResp.StatusCode)
		templates.NotFoundHandler(w, r, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	var categoriesData struct {
		Categories []domain.Category `json:"categories"`
	}

	err = helpers.DecodeBackendResponse(categoriesResp, &categoriesData)
	if err != nil {
		log.Printf("Error decoding categories response: %v", err)
		templates.NotFoundHandler(w, r, "Error loading categories", http.StatusInternalServerError)
		return
	}

	// Normalize colors
	for i := range categoriesData.Categories {
		categoriesData.Categories[i].Color = helpers.NormalizeColor(categoriesData.Categories[i].Color)
	}

	data := createPostData{
		Categories: categoriesData.Categories,
	}

	templates.RenderTemplate(w, "create_post", data)
}

// CreateTopicPost handles POST requests to /topics/create.
func (cs *ClientServer) CreateTopicPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form. File may be too large (max 20MB)", http.StatusBadRequest)
		return
	}

	allowedImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDsStr := r.Form["categories"] // This is a []string from multiple checkboxes

	// Parse category IDs
	categoryIDs := make([]int, 0)
	for _, idStr := range categoryIDsStr {
		categoryID, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			log.Printf("Invalid category ID: %v", parseErr)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	if len(categoryIDs) == 0 {
		log.Printf("No categories selected")
		http.Error(w, "At least one category must be selected", http.StatusBadRequest)
		return
	}

	// Handle optional image upload
	imagePath := ""
	file, header, err := r.FormFile("image_path")

	switch {
	case errors.Is(err, http.ErrMissingFile):
		// No image uploaded - this is fine, image is optional

	case err != nil:
		log.Printf("Error reading uploaded file: %v", err)
		http.Error(w, "Error processing uploaded file", http.StatusBadRequest)
		return

	default:
		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		if !allowedImageTypes[contentType] {
			log.Printf("Invalid file type: %s", contentType)
			http.Error(w, "Invalid file type. Only JPEG, PNG, and GIF are allowed", http.StatusBadRequest)
			return
		}

		if header.Size > maxUploadSize {
			log.Printf("File too large: %d bytes", header.Size)
			http.Error(w, "File too large. Maximum size is 20MB", http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(header.Filename)
		uniqueFilename := uuid.New().String() + ext

		err = os.MkdirAll(uploadDir, uploadDirPerm)
		if err != nil {
			log.Printf("Failed to create upload directory: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		destPath := filepath.Join(uploadDir, uniqueFilename)
		destPath = filepath.Clean(destPath)

		if !strings.HasPrefix(destPath, filepath.Clean(uploadDir)+string(os.PathSeparator)) {
			log.Printf("Invalid file path: %s", destPath)
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		var destFile *os.File
		destFile, err = os.Create(destPath)
		if err != nil {
			log.Printf("Failed to create destination file: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			log.Printf("Failed to save image: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		imagePath = "/static/images/uploads/" + uniqueFilename
	}

	createRequest := &createTopicRequest{
		CategoryIDs: categoryIDs,
		Title:       title,
		Content:     content,
		ImagePath:   imagePath,
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := cs.newRequestWithCookies(ctx, http.MethodPost, backendCreateTopic(), createRequest, r)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		// If image was uploaded, clean it up since topic creation failed
		if imagePath != "" {
			cleanupImage(imagePath)
		}
		templates.NotFoundHandler(w, r, "Failed to create topic", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Backend returned error: %s", string(body))
		// If image was uploaded, clean it up since topic creation failed
		if imagePath != "" {
			cleanupImage(imagePath)
		}
		templates.NotFoundHandler(w, r, "Failed to create topic", resp.StatusCode)
		return
	}

	// Success! Redirect to topics list
	http.Redirect(w, r, "/topics", http.StatusSeeOther)
}

// UpdateTopicPost handles POST requests to /topics/edit.
func (cs *ClientServer) UpdateTopicPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form. File may be too large (max 20MB)", http.StatusBadRequest)
		return
	}

	allowedImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	topicIDStr := r.FormValue("topic_id")
	categoryIDsStr := r.Form["categories"]
	title := r.FormValue("title")
	content := r.FormValue("content")
	currentImagePath := r.FormValue("current_image_path")

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		log.Printf("Invalid topic ID: %v", err)
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	categoryIDs := make([]int, 0)
	for _, idStr := range categoryIDsStr {
		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			log.Printf("Invalid category ID: %v", parseErr)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		categoryIDs = append(categoryIDs, id)
	}

	if len(categoryIDs) == 0 {
		http.Error(w, "At least one category must be selected", http.StatusBadRequest)
		return
	}

	// Use current image path by default
	imagePath := currentImagePath

	file, header, err := r.FormFile("image_path")
	switch {
	case errors.Is(err, http.ErrMissingFile):
		// No new image uploaded; keep current image

	case err != nil:
		log.Printf("Error reading uploaded file: %v", err)
		http.Error(w, "Error processing uploaded file", http.StatusBadRequest)
		return

	default:
		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		if !allowedImageTypes[contentType] {
			log.Printf("Invalid file type: %s", contentType)
			http.Error(w, "Invalid file type. Only JPEG, PNG, and GIF are allowed", http.StatusBadRequest)
			return
		}

		if header.Size > maxUploadSize {
			log.Printf("File too large: %d bytes", header.Size)
			http.Error(w, "File too large. Maximum size is 20MB", http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(header.Filename)
		uniqueFilename := uuid.New().String() + ext

		err = os.MkdirAll(uploadDir, uploadDirPerm)
		if err != nil {
			log.Printf("Failed to create upload directory: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		destPath := filepath.Join(uploadDir, uniqueFilename)
		destPath = filepath.Clean(destPath)

		if !strings.HasPrefix(destPath, filepath.Clean(uploadDir)+string(os.PathSeparator)) {
			log.Printf("Invalid file path: %s", destPath)
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		var destFile *os.File
		destFile, err = os.Create(destPath)
		if err != nil {
			log.Printf("Failed to create destination file: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			log.Printf("Failed to save image: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		imagePath = "/static/images/uploads/" + uniqueFilename

		if currentImagePath != "" && currentImagePath != imagePath &&
			strings.HasPrefix(currentImagePath, "/static/images/uploads/") {
			oldFilename := strings.TrimPrefix(currentImagePath, "/static/images/uploads/")
			oldFilePath := filepath.Join(uploadDir, oldFilename)
			oldFilePath = filepath.Clean(oldFilePath)

			if strings.HasPrefix(oldFilePath, filepath.Clean(uploadDir)+string(os.PathSeparator)) {
				_ = os.Remove(oldFilePath)
			}
		}
	}

	updateRequest := &updateTopicRequest{
		TopicID:     topicID,
		CategoryIDs: categoryIDs,
		Title:       title,
		Content:     content,
		ImagePath:   imagePath,
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := cs.newRequestWithCookies(ctx, http.MethodPut, backendUpdateTopic(), updateRequest, r)
	if err != nil {
		log.Printf("Backend request failed: %v", err)
		templates.NotFoundHandler(w, r, "Failed to update topic", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Backend returned error: %s", string(body))
		templates.NotFoundHandler(w, r, "Failed to update topic", resp.StatusCode)
		return
	}

	// Redirect back to the topic page
	http.Redirect(w, r, "/topic/"+topicIDStr, http.StatusSeeOther)
}

// DeleteTopicPost handles POST requests to /topics/delete.
func (cs *ClientServer) DeleteTopicPost(w http.ResponseWriter, r *http.Request) {
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
	_, err = strconv.Atoi(topicIDStr)
	if err != nil {
		log.Printf("Invalid topic ID: %v", err)
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Fetch topic to get image path
	getURL := backendGetTopicByID() + "?id=" + topicIDStr
	getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	ip := middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(getReq, ip)

	for _, cookie := range r.Cookies() {
		getReq.AddCookie(cookie)
	}

	getResp, err := cs.HTTPClient.Do(getReq)
	if err != nil {
		http.Error(w, "Failed to fetch topic", http.StatusInternalServerError)
		return
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch topic", getResp.StatusCode)
		return
	}

	var topicResp struct {
		ImagePath string `json:"imagePath"`
	}

	err = helpers.DecodeBackendResponse(getResp, &topicResp)
	if err != nil {
		http.Error(w, "Failed to decode topic data", http.StatusInternalServerError)
		return
	}

	// Delete topic in backend
	deleteURL := backendDeleteTopic() + "?id=" + topicIDStr
	delReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		http.Error(w, "Failed to create delete request", http.StatusInternalServerError)
		return
	}

	ip = middleware.GetIPFromContext(r)
	if ip == "" {
		http.Error(w, "Error no IP found in request", http.StatusInternalServerError)
	}

	helpers.SetIPHeaders(delReq, ip)

	for _, cookie := range r.Cookies() {
		delReq.AddCookie(cookie)
	}

	delResp, err := cs.HTTPClient.Do(delReq)
	if err != nil {
		http.Error(w, "Failed to delete topic", http.StatusInternalServerError)
		return
	}
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(delResp.Body)
		log.Printf("Backend delete error: %s", string(body))
		http.Error(w, "Failed to delete topic", delResp.StatusCode)
		return
	}

	// Delete image file locally
	if topicResp.ImagePath != "" &&
		strings.HasPrefix(topicResp.ImagePath, "/static/images/uploads/") {
		filename := strings.TrimPrefix(topicResp.ImagePath, "/static/images/uploads/")
		filePath := filepath.Join(uploadDir, filename)

		err = os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to delete image file %s: %v", filePath, err)
		}
	}

	http.Redirect(w, r, "/topics", http.StatusSeeOther)
}

// Helper function to clean up uploaded image if topic creation fails.
func cleanupImage(imagePath string) {
	if imagePath != "" && strings.HasPrefix(imagePath, "/static/images/uploads/") {
		filename := strings.TrimPrefix(imagePath, "/static/images/uploads/")
		filePath := filepath.Join(uploadDir, filename)
		filePath = filepath.Clean(filePath)

		if strings.HasPrefix(filePath, filepath.Clean(uploadDir)+string(os.PathSeparator)) {
			_ = os.Remove(filePath)
		}
	}
}
