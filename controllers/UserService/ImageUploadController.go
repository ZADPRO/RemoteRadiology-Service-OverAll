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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostUploadProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()
		log.Println("\n\nIncoming request to upload profile image\n")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			log.Println("\n\nMissing user ID or roleId in request context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID not found in request context.",
			})
			return
		}
		log.Printf("User ID: %v, Role ID: %v\n", idValue, roleIdValue)

		file, err := c.FormFile("profileImage")
		if err != nil {
			log.Printf("Error retrieving file: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Error retrieving profile image: " + err.Error(),
			})
			return
		}
		log.Printf("\nReceived file: %s (%d bytes)\n", file.Filename, file.Size)

		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			log.Printf("Invalid file type: %s\n", ext)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid profile image type. Only JPG, JPEG, PNG allowed.",
			})
			return
		}
		log.Printf("File extension validated: %s\n", ext)

		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),
			timeZone.GetTimeWithFormate("20060102150405"),
			ext)
		log.Printf("Unique filename generated: %s\n", uniqueFilename)

		s3Key := s3path.BuildS3Key(file.Filename, uniqueFilename)
		log.Printf("S3 key built: %s\n", s3Key)

		uploadURL, err := s3Service.GeneratePresignPutURL(c, s3Key, 15*60*1e9)
		if err != nil {
			log.Printf("Error generating presigned URL for %s: %v\n", s3Key, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate S3 upload URL",
			})
			return
		}
		log.Printf("Presigned upload URL generated for key: %s\n", s3Key)

		payload := map[string]interface{}{
			"status":    true,
			"message":   "Presigned URL generated successfully. Upload file to this URL.",
			"fileName":  uniqueFilename,
			"s3Key":     s3Key,
			"uploadURL": uploadURL,
		}
		token := accesstoken.CreateToken(idValue, roleIdValue)

		log.Println("\n\nProfile image upload presigned URL response sent successfully", payload)
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostUploadPublicProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()
		log.Println("Incoming request for public profile image upload")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Unauthorized: User ID or roleId missing",
			})
			return
		}

		var req struct {
			Extension string `json:"extension"` // send only file type from frontend, e.g., ".jpg"
		}

		if err := c.BindJSON(&req); err != nil || req.Extension == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request. Extension is required.",
			})
			return
		}

		ext := strings.ToLower(req.Extension)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid file type. Only JPG, JPEG, PNG allowed.",
			})
			return
		}

		// Generate unique filename
		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),
			timeZone.GetTimeWithFormate("20060102150405"),
			ext)

		s3Key := fmt.Sprintf("images/%s", uniqueFilename)

		// Generate presigned URL for frontend upload
		uploadURL, err := s3Service.GeneratePresignPutURLPublic(c, s3Key, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate S3 presigned URL",
			})
			return
		}

		viewURL := fmt.Sprintf("https://%s.s3.us-east-2.amazonaws.com/%s", s3Service.GetPublicBucketName(), s3Key)

		payload := map[string]interface{}{
			"status":    true,
			"message":   "Public presigned URL generated successfully",
			"fileName":  uniqueFilename,
			"s3Key":     s3Key,
			"uploadURL": uploadURL,
			"viewURL":   viewURL,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostUploadPrivateDocument() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.InitLogger()
		log.Println("Incoming request for private document upload")

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Unauthorized: User ID or roleId missing",
			})
			return
		}

		var req struct {
			Extension string `json:"extension"` // e.g. ".pdf", ".docx"
		}
		if err := c.BindJSON(&req); err != nil || req.Extension == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request. Extension is required.",
			})
			return
		}

		ext := strings.ToLower(req.Extension)
		allowed := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"}
		valid := false
		for _, a := range allowed {
			if ext == a {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid file type. Allowed: PDF, DOC, DOCX, XLS, XLSX, TXT.",
			})
			return
		}

		// Create unique file name
		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),
			timeZone.GetTimeWithFormate("20060102150405"),
			ext,
		)

		s3Key := fmt.Sprintf("documents/%s", uniqueFilename)

		// Generate presigned PUT URL for upload
		uploadURL, err := s3Service.GeneratePresignPutURLPrivate(c, s3Key, 15*time.Minute)
		if err != nil {
			log.Errorf("Error generating presigned PUT URL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate upload URL",
			})
			return
		}

		// Generate presigned GET (view/download) URL - valid for 10 minutes
		viewURL, err := s3Service.GeneratePresignGetURLPrivate(c, s3Key, 10*time.Minute)
		if err != nil {
			log.Errorf("Error generating presigned GET URL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate view/download URL",
			})
			return
		}

		payload := map[string]interface{}{
			"status":      true,
			"message":     "Private presigned URLs generated successfully",
			"fileName":    uniqueFilename,
			"s3Key":       s3Key,
			"uploadURL":   uploadURL,
			"viewURL":     viewURL,
			"expiresIn":   "10 minutes",
			"accessLevel": "private",
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
