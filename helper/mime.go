package helper

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// GetMimeType detects MIME type from reader and returns buffer, content type, extension, and error
func GetMimeType(reader io.Reader) (*bytes.Buffer, string, string, error) {
	// Read first 512 bytes to detect content type
	peek := make([]byte, 512)
	n, err := reader.Read(peek)
	if err != nil && err != io.EOF {
		return nil, "", "", err
	}
	peek = peek[:n]

	// Detect content type
	contentType := http.DetectContentType(peek)

	// Get extension from content type
	extension := GetExtensionFromMimeType(contentType)

	// Create buffer with all content (peek + remaining)
	buffer := new(bytes.Buffer)
	buffer.Write(peek)

	// Read remaining data from reader
	if remaining, err := io.ReadAll(reader); err == nil {
		buffer.Write(remaining)
	}

	return buffer, contentType, extension, nil
}

// GetExtensionFromMimeType returns file extension from MIME type
func GetExtensionFromMimeType(mimeType string) string {
	// Remove charset if present
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = mimeType[:idx]
	}
	mimeType = strings.TrimSpace(mimeType)

	mimeToExt := map[string]string{
		// Images
		"image/jpeg":    ".jpg",
		"image/png":     ".png",
		"image/gif":     ".gif",
		"image/webp":    ".webp",
		"image/svg+xml": ".svg",
		"image/bmp":     ".bmp",
		"image/tiff":    ".tiff",

		// Videos
		"video/mp4":        ".mp4",
		"video/mpeg":       ".mpeg",
		"video/quicktime":  ".mov",
		"video/x-msvideo":  ".avi",
		"video/x-matroska": ".mkv",
		"video/webm":       ".webm",

		// Audio
		"audio/mpeg": ".mp3",
		"audio/wav":  ".wav",
		"audio/ogg":  ".ogg",

		// Documents
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"text/plain":       ".txt",
		"text/html":        ".html",
		"application/json": ".json",

		// Archives
		"application/zip":             ".zip",
		"application/x-rar":           ".rar",
		"application/x-7z-compressed": ".7z",
	}

	if ext, ok := mimeToExt[mimeType]; ok {
		return ext
	}

	return ""
}

// GetFileExtension returns file extension from filename
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsImageFile checks if file is an image based on filename
func IsImageFile(filename string) bool {
	ext := GetFileExtension(filename)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".tiff"}
	for _, e := range imageExts {
		if ext == e {
			return true
		}
	}
	return false
}

// IsVideoFile checks if file is a video based on filename
func IsVideoFile(filename string) bool {
	ext := GetFileExtension(filename)
	videoExts := []string{".mp4", ".avi", ".mov", ".mkv", ".webm", ".mpeg"}
	for _, e := range videoExts {
		if ext == e {
			return true
		}
	}
	return false
}

// ReadFullBuffer reads the remaining data and combines with initial buffer
func ReadFullBuffer(initialBuffer []byte, reader io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(initialBuffer)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
