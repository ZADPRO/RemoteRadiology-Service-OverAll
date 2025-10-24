package s3path

import (
	"path/filepath"
	"strings"
)

func GetS3Folder(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff":
		return "images"
	case ".pdf", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".txt":
		return "documents"
	case ".dcm":
		return "dicom"
	default:
		return "others"
	}
}

func BuildS3Key(filename, uniqueName string) string {
	folder := GetS3Folder(filename)
	return folder + "/" + uniqueName
}

func BuildFinalReportKey(filename string) string {
	return "finalReport/" + filename
}
