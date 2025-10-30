package s3path

import (
	s3Service "AuthenticationService/Service/S3"
	"context"

)

func GetS3FileURL(fileName string) (string, error) {
	if fileName == "" {
		return "", nil
	}

	ctx := context.Background()
	return s3Service.GeneratePresignGetURL(ctx, fileName, 15*60*1e9)
}
