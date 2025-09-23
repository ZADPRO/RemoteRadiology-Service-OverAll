package s3config

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *S3Client) PresignGet(ctx context.Context, key string, expire time.Duration) (string, error) {
	output, err := s.Presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = expire
	})

	if err != nil {
		return "", err
	}

	return output.URL, nil
}

func (s *S3Client) PresignPut(ctx context.Context, key string, expire time.Duration) (string, error) {
	output, err := s.Presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = expire
	})

	if err != nil {
		return "", err
	}

	return output.URL, nil

}
