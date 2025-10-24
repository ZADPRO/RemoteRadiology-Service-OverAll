package s3config

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *S3Client) UploadFromReader(ctx context.Context, key string, body io.Reader, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}

	_, err := s.Client.PutObject(ctx, input)
	return err
}
