package s3Controller

import (
	s3Service "AuthenticationService/Service/S3"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

type UploadFilePayload struct {
	FileType      string `json:"fileType"` // "consentForm" or "finalReport"
	FileUrl       string `json:"fileUrl"`
	PatientId     int    `json:"patientId"`
	AppointmentId int    `json:"appointmentId"`
}

func S3FinalReportUploadController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from JWT context
		idValue, idExists := c.Get("id")
		roleIdValue, _ := c.Get("roleId")

		if !idExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID not found in request context.",
			})
			return
		}

		// Parse request body
		var payload UploadFilePayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid payload",
				"error":   err.Error(),
			})
			return
		}

		if payload.FileType != "consentForm" && payload.FileType != "finalReport" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid fileType",
			})
			return
		}

		// DB connection
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// Prepare values
		consentForm := ""
		finalReport := ""
		if payload.FileType == "consentForm" {
			consentForm = payload.FileUrl
		} else {
			finalReport = payload.FileUrl
		}

		// 1️⃣ Check if user exists
		var existingUser struct {
			ID int
		}
		userCheck := dbConn.Table(`"backupFiles".report`).Where(`"userId" = ?`, payload.PatientId).First(&existingUser)

		if userCheck.Error != nil && userCheck.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Database query failed",
				"error":   userCheck.Error.Error(),
			})
			return
		}

		if userCheck.RowsAffected > 0 {
			// 2️⃣ Check if record exists for this user + appointment
			var report struct {
				ID int
			}
			recordCheck := dbConn.Table(`"backupFiles".report`).
				Where(`"userId" = ? AND "appointmentId" = ?`, payload.PatientId, payload.AppointmentId).
				First(&report)

			if recordCheck.Error != nil && recordCheck.Error != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Database query failed",
					"error":   recordCheck.Error.Error(),
				})
				return
			}

			if recordCheck.RowsAffected > 0 {
				// ✅ Update existing record for this appointment
				updates := map[string]interface{}{}
				if payload.FileType == "consentForm" {
					updates["consentForm"] = consentForm
				} else {
					updates["finalReportPath"] = finalReport
				}

				if err := dbConn.Table(`"backupFiles".report`).
					Where(`"userId" = ? AND "appointmentId" = ?`, payload.PatientId, payload.AppointmentId).
					Updates(updates).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  false,
						"message": "Failed to update record",
						"error":   err.Error(),
					})
					return
				}
			} else {
				// 🆕 User exists, but this is a new appointment -> insert new record
				newReport := map[string]interface{}{
					"userId":          payload.PatientId,
					"appointmentId":   payload.AppointmentId,
					"consentForm":     consentForm,
					"finalReportPath": finalReport,
				}
				if err := dbConn.Table(`"backupFiles".report`).Create(&newReport).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  false,
						"message": "Failed to insert record",
						"error":   err.Error(),
					})
					return
				}
			}
		} else {
			// 🆕 User does not exist -> insert new record
			newReport := map[string]interface{}{
				"userId":          payload.PatientId,
				"appointmentId":   payload.AppointmentId,
				"consentForm":     consentForm,
				"finalReportPath": finalReport,
			}
			if err := dbConn.Table(`"backupFiles".report`).Create(&newReport).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to insert record",
					"error":   err.Error(),
				})
				return
			}
		}

		// Generate presigned URL for frontend
		s3URL, err := s3Service.GenerateFinalReportPresignURL(c, payload.FileUrl, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate presigned URL",
				"error":   err.Error(),
			})
			return
		}

		// Create new JWT token
		token := accesstoken.CreateToken(idValue, roleIdValue)

		respPayload := map[string]interface{}{
			"status":  true,
			"message": "File info saved successfully",
			"url":     s3URL,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  respPayload,
			"token": token,
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
