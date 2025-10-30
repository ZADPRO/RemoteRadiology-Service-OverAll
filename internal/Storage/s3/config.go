package s3config

import (
	logger "AuthenticationService/internal/Helper/Logger"
	"os"

)

type S3Config struct {
	Region string
	Bucket string
}

func LoadConfig() *S3Config {
	logger := logger.InitLogger()

	region := os.Getenv("AWS_REGION")
	bucket := "easeqt-health-archive"

	if region == "" || bucket == "" {
		logger.Info("❌ Missing AWS_REGION or AWS_S3_BUCKET in environment")
	}

	return &S3Config{
		Region: region,
		Bucket: bucket,
	}
}
