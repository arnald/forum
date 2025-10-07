// Package handlers - upload.go implements secure file upload functionality.
// This file handles image uploads for posts, including validation of file types,
// size limits, and secure filename generation. Implements proper security measures
// to prevent malicious uploads and directory traversal attacks.
package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// Maximum file size for uploads (20MB)
// This prevents DoS attacks via large file uploads and manages disk space
const maxUploadSize = 20 * 1024 * 1024 // 20MB

// allowedExtensions defines which file extensions are permitted for upload
// Using a whitelist approach for security - only explicitly allowed types
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// allowedMimeTypes defines which MIME types are permitted for upload
// This provides an additional layer of validation beyond file extensions
// to prevent malicious files disguised with valid extensions
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// handleImageUpload processes file uploads with comprehensive security validation.
// Performs multiple layers of security checks:
// 1. File size validation to prevent DoS attacks
// 2. Extension validation to block unwanted file types
// 3. MIME type validation to detect disguised malicious files
// 4. Secure filename generation to prevent path traversal attacks
// Returns the web-accessible path to the uploaded file or an error
func (h *Handler) handleImageUpload(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size to prevent large uploads
	if header.Size > maxUploadSize {
		return "", fmt.Errorf("file too large: %d bytes (max %d bytes)", header.Size, maxUploadSize)
	}

	// Validate file extension (first layer of defense)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	// Read first 512 bytes to detect actual file type (second layer of defense)
	// This prevents malicious files disguised with valid extensions
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Validate MIME type based on file content (not just extension)
	// Uses Go's built-in content type detection
	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return "", fmt.Errorf("invalid file type: %s", mimeType)
	}

	// Reset file pointer to beginning for subsequent reading
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	// Generate unique filename using UUID to prevent:
	// - Filename collisions
	// - Path traversal attacks
	// - Overwriting existing files
	imageID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	// Ensure upload directory exists
	uploadDir := "static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Construct secure filename: UUID_timestamp_extension
	// Timestamp adds additional uniqueness and helps with debugging
	filename := fmt.Sprintf("%s_%d%s", imageID.String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy uploaded file content to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	// Return web-accessible path (with leading slash for URL)
	return "/static/uploads/" + filename, nil
}
