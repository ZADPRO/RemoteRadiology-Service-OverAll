package helper

import (
	"encoding/base64"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
)

type FileData struct {
	Base64Data  string `json:"base64Data"`  // base64-encoded file content
	ContentType string `json:"contentType"` // e.g., "image/jpeg"
}

func ViewFile(filePath string) (*FileData, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Convert file bytes to base64
	base64Content := base64.StdEncoding.EncodeToString(data)

	// Get MIME type
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream" // fallback
	}

	return &FileData{
		Base64Data:  base64Content,
		ContentType: contentType,
	}, nil
}
