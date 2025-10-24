package s3Controller

import (
	s3Service "AuthenticationService/Service/S3"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func S3GeneratePresignPutController() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.JSON(400, gin.H{"status": false, "message": "Missing filename"})
			return
		}

		url, err := s3Service.GeneratePresignPutURL(c, filename, 15*time.Minute)
		if err != nil {
			c.JSON(500, gin.H{"status": false, "message": "Failed to generate presign URL"})
			return
		}

		c.JSON(200, gin.H{"status": true, "url": url})
	}
}

func S3GeneratePresignGetController() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Missing key parameter"})
			return
		}

		url, err := s3Service.GeneratePresignGetURL(c, key, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to generate download URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Download URL generated successfully",
			"url":     url,
		})
	}
}

func S3GetFileController() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(400, gin.H{"status": false, "message": "Missing key parameter"})
			return
		}

		url, err := s3Service.GeneratePresignURL(key, 15) // 15 minutes expiry
		if err != nil {
			c.JSON(500, gin.H{"status": false, "message": "Failed to generate presigned URL", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"status":  true,
			"message": "Presigned URL generated successfully",
			"url":     url,
		})
	}
}

func AckCheckController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Always return true for now
		c.JSON(http.StatusOK, gin.H{
			"status": true,
		})
	}
}

func S3FinalReportUploadController() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Missing filename parameter",
			})
			return
		}

		url, err := s3Service.GenerateFinalReportPresignURL(c, filename, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate presigned URL",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Presigned upload URL generated successfully",
			"url":     url,
		})
	}
}

func S3PublicProfileUploadController() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Missing filename parameter",
			})
			return
		}

		url, err := s3Service.GeneratePresignPutURLPublic(c, filename, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate presigned URL",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Presigned upload URL generated successfully",
			"url":     url,
		})
	}
}
