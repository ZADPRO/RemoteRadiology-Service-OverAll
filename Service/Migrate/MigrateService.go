package service

import (
	logger "AuthenticationService/internal/Helper/Logger"
	model "AuthenticationService/internal/Model/Migrate"
	query "AuthenticationService/query/Migrate"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
)

func UploadFileToS3(bucketName, region, keyName, localFilePath string) (string, error) {
	// Load AWS credentials and config from environment or shared credentials
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Open local file
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Upload input
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	}

	// Upload to S3
	_, err = client.PutObject(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Construct S3 URL (if bucket is public or uses standard URL)
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, keyName)

	return s3URL, nil
}

func DicomMigrateService(db *gorm.DB) (bool, string) {
	log := logger.InitLogger()

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("ERROR: Failed to begin transaction: %v\n", tx.Error)
		return false, "Something went wrong, Try Again"
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("ERROR: Recovered from panic, rolling back transaction:", r)
			tx.Rollback()
		}
	}()

	MigrationDicomFile := []model.MigrationDicomFile{}

	CheckMigrateErr := tx.Raw(query.CheckMigrateDicomSQL).Scan(&MigrationDicomFile).Error
	if CheckMigrateErr != nil {
		log.Printf("ERROR: Failed to execute query: %v\n", CheckMigrateErr)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	if len(MigrationDicomFile) > 0 {

		for _, v := range MigrationDicomFile {

			if !v.IsMigrated {

				localPath := filepath.Join("Assets/Dicom/", v.RefDFFilename)
				key := "dicom/" + v.RefDFFilename // path inside your S3 bucket
				bucket := os.Getenv("AWS_S3_BUCKET")
				region := os.Getenv("AWS_REGION")

				url, err := UploadFileToS3(bucket, region, key, localPath)
				if err != nil {
					log.Printf("upload failed: %v", err)
				}

				MirgateDicomNewErr := tx.Exec(
					query.NewMigrateDicomSQL,
					v.RefDFId,
					v.RefUserId,
					v.RefAppointmentId,
					v.RefDFFilename,
					url,
				).Error
				if MirgateDicomNewErr != nil {
					log.Printf("ERROR: Failed to execute query: %v\n", MirgateDicomNewErr)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

				UpdateDicomErr := tx.Exec(
					query.UpdateDicomSQL,
					url,
					v.RefDFId,
				).Error
				if UpdateDicomErr != nil {
					log.Printf("ERROR: Failed to execute query: %v\n", UpdateDicomErr)
					tx.Rollback()
					return false, "Something went wrong, Try Again"
				}

				fmt.Println("File Name: ", v.RefDFFilename)
			}

		}

	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Failed to commit transaction: %v\n", err)
		tx.Rollback()
		return false, "Something went wrong, Try Again"
	}

	return true, "Succcessfully Dicom Migrated"
}
