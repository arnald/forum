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

const maxUploadSize = 20 * 1024 * 1024 // 20MB

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func (h *Handler) handleImageUpload(file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > maxUploadSize {
		return "", fmt.Errorf("file too large: %d bytes (max %d bytes)", header.Size, maxUploadSize)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return "", fmt.Errorf("invalid file type: %s", mimeType)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	imageID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	uploadDir := "static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s_%d%s", imageID.String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return "/static/uploads/" + filename, nil
}
