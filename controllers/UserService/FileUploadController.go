package controllers

import (
	s3Service "AuthenticationService/Service/S3"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	logger "AuthenticationService/internal/Helper/Logger"
	s3path "AuthenticationService/internal/Helper/S3"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostUploadFileController() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			log.Println("User ID or RoleID missing in request context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID not found in request context.",
			})
			return
		}
		log.Printf("Request received for userID: %v, roleID: %v\n", idValue, roleIdValue)

		file, err := c.FormFile("file")
		if err != nil {
			log.Printf("Error retrieving file from request: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Error retrieving file: " + err.Error(),
			})
			return
		}
		log.Printf("File received: %s, size: %d bytes\n", file.Filename, file.Size)

		const maxFileSize = 10 * 1024 * 1024
		if file.Size > maxFileSize {
			log.Printf("File size exceeds limit: %d bytes\n", file.Size)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": fmt.Sprintf("File size exceeds %d MB limit", maxFileSize/(1024*1024)),
			})
			return
		}

		ext := filepath.Ext(file.Filename)
		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),
			timeZone.GetTimeWithFormate("20060102150405"),
			ext,
		)
		log.Printf("Generated unique filename: %s\n", uniqueFilename)

		s3Key := s3path.BuildS3Key(file.Filename, uniqueFilename)
		log.Printf("S3 key for upload: %s\n", s3Key)

		uploadURL, err := s3Service.GeneratePresignPutURL(c, s3Key, 15*60*1e9)
		if err != nil {
			log.Printf("Error generating presigned URL for S3 upload: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate S3 upload URL",
			})
			return
		}
		log.Printf("Presigned S3 upload URL generated successfully\n")

		payload := map[string]interface{}{
			"status":      true,
			"message":     "Presigned URL generated successfully. Upload file to this URL.",
			"fileName":    uniqueFilename,
			"oldFileName": file.Filename,
			"s3Key":       s3Key,
			"uploadURL":   uploadURL,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)
		log.Printf("Token generated for userID: %v\n", idValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
		log.Println("Response sent successfully for file upload request")
	}
}
