package s3Controller

import (
	s3Service "AuthenticationService/Service/S3"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	"context"
	"fmt"
	"net/http"
	"strings"
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

type S3DailyBackupResponse struct {
	ScanCenterName string                         `json:"ScanCenterName"`
	Users          map[string]map[string][]string `json:"Users"` // UserCustId -> FileType -> []Files
}

func S3DailyBackupController() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		type row struct {
			ScanCenterId      int
			ScanCenterName    string
			UserId            int
			UserCustId        string
			DicomFileName     string
			AppointmentId     int
			OldReportCategory int
			OldReportFileName string
			ConsentFormPath   string
			FinalReportPath   string
		}

		var rows []row
		query := `
		SELECT
    sc."refSCId" AS "ScanCenterId",
    sc."refSCName" AS "ScanCenterName",
    u."refUserId" AS "UserId",
    u."refUserCustId" AS "UserCustId",
    df."refDFFilename" AS "DicomFileName",
    df."refAppointmentId" AS "AppointmentId",
    orp."refORCategoryId" AS "OldReportCategoryId",
    orp."refORFilename" AS "OldReportFileName",
    r."consentForm" AS "ConsentFormPath",
    r."finalReportPath" AS "FinalReportPath"
FROM
    public."ScanCenter" sc
LEFT JOIN
    map."refScanCenterMapPatient" rscmp
    ON sc."refSCId" = rscmp."refSCId"
LEFT JOIN
    public."Users" u
    ON rscmp."refUserId" = u."refUserId"
LEFT JOIN
    dicom."refDicomFiles" df
    ON u."refUserId" = df."refUserId"
LEFT JOIN
    notes."refOldReport" orp
    ON u."refUserId" = orp."refUserId"
LEFT JOIN
    "backupFiles".report r
    ON u."refUserId" = r."userId"
ORDER BY
    sc."refSCId",
    u."refUserId",
    df."refAppointmentId",
    orp."refORCategoryId";` 

		if err := dbConn.Raw(query).Scan(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "DB query failed", "error": err.Error()})
			return
		}

		responseMap := make(map[string]*S3DailyBackupResponse)

		for _, r := range rows {
			scanCenterName := hashdb.Decrypt(r.ScanCenterName)
			userCustId := hashdb.Decrypt(r.UserCustId)

			if _, ok := responseMap[scanCenterName]; !ok {
				responseMap[scanCenterName] = &S3DailyBackupResponse{
					ScanCenterName: scanCenterName,
					Users:          make(map[string]map[string][]string),
				}
			}

			userFiles := responseMap[scanCenterName].Users
			if _, ok := userFiles[userCustId]; !ok {
				userFiles[userCustId] = map[string][]string{
					"DicomFiles":  {},
					"OldReports":  {},
					"ConsentForm": {},
					"FinalReport": {},
				}
			}

			// Function to generate presigned URL for any file
			getPresignURL := func(fileName, folder string) string {
				if fileName == "" {
					return ""
				}

				// If the fileName is already a full URL, extract the last segment as key
				key := fileName
				if strings.HasPrefix(fileName, "https://") {
					parts := strings.Split(fileName, "/")
					key = fmt.Sprintf("%s/%s", folder, parts[len(parts)-1])
				} else {
					key = fmt.Sprintf("%s/%s", folder, fileName)
				}

				url, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Hour)
				if err != nil {
					fmt.Println("Error generating presign for:", key, err)
					return ""
				}
				return url
			}

			// Generate presigned URLs for all files
			if r.DicomFileName != "" {
				userFiles[userCustId]["DicomFiles"] = append(userFiles[userCustId]["DicomFiles"], getPresignURL(r.DicomFileName, "dicom"))
			}
			if r.OldReportFileName != "" {
				userFiles[userCustId]["OldReports"] = append(userFiles[userCustId]["OldReports"], getPresignURL(r.OldReportFileName, "oldReportsPatient"))
			}
			if r.ConsentFormPath != "" {
				userFiles[userCustId]["ConsentForm"] = append(userFiles[userCustId]["ConsentForm"], getPresignURL(r.ConsentFormPath, "finalReport"))
			}
			if r.FinalReportPath != "" {
				userFiles[userCustId]["FinalReport"] = append(userFiles[userCustId]["FinalReport"], getPresignURL(r.FinalReportPath, "finalReport"))
			}
		}

		var response []S3DailyBackupResponse
		for _, v := range responseMap {
			response = append(response, *v)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   response,
		})
	}
}
