package s3config

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	Client     *s3.Client
	Uploader   *manager.Uploader
	Downloader *manager.Downloader
	Presign    *s3.PresignClient
	Bucket     string
}

func New(ctx context.Context) (*S3Client, error) {
	appCfg := LoadConfig()

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		Client:     client,
		Uploader:   manager.NewUploader(client),
		Downloader: manager.NewDownloader(client),
		Presign:    s3.NewPresignClient(client),
		Bucket:     appCfg.Bucket,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, key string, body io.Reader) error {
	_, err := s.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   body,
		ACL:    types.ObjectCannedACLPublicRead, // ✅ Correct type
	})
	return err
}
func (s *S3Client) DownloadBytes(ctx context.Context, fileURL string) ([]byte, error) {
	var key string

	// If the input is a full S3 URL (https://bucket.s3.region.amazonaws.com/key)
	if strings.HasPrefix(fileURL, "http") {
		u, err := url.Parse(fileURL)
		if err != nil {
			return nil, err
		}
		// The key is the path without the leading slash
		key = strings.TrimPrefix(u.Path, "/")
	} else {
		// Otherwise, treat as key directly
		key = fileURL
	}

	if key == "" {
		return nil, errors.New("S3 key is empty")
	}

	// Use s3manager.Downloader to download the file into a buffer
	buff := manager.NewWriteAtBuffer([]byte{})

	_, err := s.Downloader.Download(ctx, buff, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
