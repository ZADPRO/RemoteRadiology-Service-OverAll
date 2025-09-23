package s3Service

import (
	s3config "AuthenticationService/internal/Storage/s3"
	"context"
	"time"
)

var s3Client *s3config.S3Client

func init() {
	client, err := s3config.New(context.Background())
	if err != nil {
		panic(err)
	}
	s3Client = client
}

func GeneratePresignPutURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	return s3Client.PresignPut(ctx, key, expire)
}

func GeneratePresignGetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
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
