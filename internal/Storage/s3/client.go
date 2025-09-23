package s3config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
