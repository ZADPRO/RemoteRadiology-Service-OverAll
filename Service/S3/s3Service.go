package s3Service

import (
	s3config "AuthenticationService/internal/Storage/s3"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

)

func getS3Folder(filename string) string {
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

var s3Client *s3config.S3Client

func init() {
	client, err := s3config.New(context.Background())
	if err != nil {
		panic(err)
	}
	s3Client = client
}

func GeneratePresignPutURL(ctx context.Context, filename string, expire time.Duration) (string, error) {
	folder := getS3Folder(filename)
	key := folder + "/" + filename

	return s3Client.PresignPut(ctx, key, expire)
}

func GeneratePresignGetURL(ctx context.Context, filename string, expire time.Duration) (string, error) {
	// folder := getS3Folder(filename)
	key := filename
	fmt.Printf("\n\n\nFolder %v \n\nFile name %v\n\n", key, filename)

	return s3Client.PresignGet(ctx, key, expire)
}

func GeneratePresignURL(key string, expireMinutes int) (string, error) {
	ctx := context.Background()
	expire := time.Duration(expireMinutes) * time.Minute
	url, err := s3Client.PresignGet(ctx, key, expire)
	if err != nil {
		return "", err
	}
	return url, nil
}

func UploadFinalReportToS3(ctx context.Context, filename string, file io.Reader) (string, error) {
	key := fmt.Sprintf("finalReport/%s", filename)
	err := s3Client.UploadFile(ctx, key, file)
	if err != nil {
		return "", err
	}

	// Optionally generate a public presigned GET URL after upload
	url, err := s3Client.PresignGet(ctx, key, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}

func GenerateFinalReportPresignURL(ctx context.Context, filename string, expire time.Duration) (string, error) {
	key := fmt.Sprintf("finalReport/%s", filename)
	return s3Client.PresignPut(ctx, key, expire)
}

func GeneratePresignPutURLPublic(ctx context.Context, key string, expire time.Duration) (string, error) {
	return s3Client.PresignPutPublic(ctx, key, expire)
}

func GetPublicBucketName() string {
	return s3Client.Bucket
}

// GeneratePresignPutURLPrivate creates presigned PUT URL for private document upload
func GeneratePresignPutURLPrivate(ctx context.Context, key string, expire time.Duration) (string, error) {
	return s3Client.PresignPut(ctx, key, expire) // private upload (no public ACL)
}

// GeneratePresignGetURLPrivate creates temporary presigned GET URL for private document download
func GeneratePresignGetURLPrivate(ctx context.Context, key string, expire time.Duration) (string, error) {
	return s3Client.PresignGet(ctx, key, expire)
}
