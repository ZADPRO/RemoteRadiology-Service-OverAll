package controllers

import (
	service "AuthenticationService/Service/Appointment"
	s3Service "AuthenticationService/Service/S3"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	model "AuthenticationService/internal/Model/Appointment"
	s3config "AuthenticationService/internal/Storage/s3"
	query "AuthenticationService/query/Appointment"
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddTechnicianIntakeFormController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		fmt.Println("Allowed")

		data, ok := helper.RequestHandler[model.AddTechnicianIntakeFormReq](c)
		if !ok {
			return
		}

		// var reqVal model.AddIntakeFormReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AddTechnicianIntakeFormService(dbConn, *data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func ViewTechnicianIntakeFormController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		fmt.Println("Allowed")

		data, ok := helper.RequestHandler[model.ViewTechnicianIntakeFormReq](c)
		if !ok {
			return
		}

		// var reqVal model.AddIntakeFormReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		ViewIntakeData, Aduit, TechIntakeData, TechnicianName, TechnicianCustId := service.ViewTechnicianIntakeFormService(dbConn, *data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"IntakeData":       ViewIntakeData,
			"IntakeDataAduit":  Aduit,
			"TechIntakeData":   TechIntakeData,
			"technicianName":   TechnicianName,
			"technicianCustId": TechnicianCustId,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AssignTechnicianController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.ViewTechnicianIntakeFormReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		Status, Message := service.AssignTechnicianService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  Status,
			"message": Message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostUploadDicomFileController() gin.HandlerFunc {
	return func(c *gin.Context) {

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID not found in request context.",
			})
			return
		}

		log := logger.InitLogger()

		// Retrieve file from multipart form
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Error retrieving file: " + err.Error(),
			})
			return
		}

		maxFileSize := int64(5500 * 1024 * 1024) // 5.5 GB
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": fmt.Sprintf("File size exceeds the limit of %d MB", maxFileSize/(1024*1024)),
			})
			return
		}

		ext := filepath.Ext(file.Filename)
		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),
			timeZone.GetTimeWithFormate("20060102150405"),
			ext,
		)

		// S3 Key: store under /dicom/
		s3Key := fmt.Sprintf("dicom/%s", uniqueFilename)

		// Open the uploaded file
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to open uploaded file",
			})
			return
		}
		defer srcFile.Close()

		// Initialize S3 client
		ctx := c.Request.Context()
		s3Client, err := s3config.New(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to initialize S3 client",
			})
			return
		}

		// Upload file to S3
		if err := s3Client.UploadFromReader(ctx, s3Key, srcFile, file.Header.Get("Content-Type")); err != nil {
			log.Printf("Error uploading file to S3: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to upload file to S3",
			})
			return
		}

		// Generate a presigned GET URL (optional) for temporary access
		viewURL, err := s3Client.PresignGet(ctx, s3Key, 15*time.Minute)
		if err != nil {
			log.Printf("Error generating presigned GET URL: %v", err)
			viewURL = "" // optional, fallback
		}

		payload := map[string]interface{}{
			"status":      true,
			"message":     "Dicom file uploaded successfully!",
			"fileName":    uniqueFilename,
			"s3Key":       s3Key,
			"viewURL":     viewURL,
			"oldFilename": file.Filename,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func GenerateDicomFileName(scanCenterCustId, patientCustId, side, originalFileName string) string {
	ext := filepath.Ext(originalFileName)
	if ext == "" {
		ext = ".zip"
	}

	currentDate := timeZone.GetTimeWithFormate("02-01-2006")
	timestamp := time.Now().UnixMilli()

	sideCode := "R"
	if strings.ToLower(side) == "left" {
		sideCode = "L"
	}

	return fmt.Sprintf("%s_%s_%s_%s_%d%s",
		scanCenterCustId,
		strings.ToUpper(patientCustId),
		currentDate,
		sideCode,
		timestamp,
		ext,
	)
}

func PostGenerateDicomUploadURLController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Step 1: Validate user context ---
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID or Role ID not found in context.",
			})
			return
		}

		// --- Step 2: Parse request ---
		var req struct {
			FileName      string `json:"fileName"`
			Side          string `json:"side"`
			AppointmentId int    `json:"appointmentId"`
			PatientId     int    `json:"patientId"`
		}

		if err := c.BindJSON(&req); err != nil || req.FileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Missing or invalid request payload.",
			})
			return
		}

		// --- Step 3: Initialize DB ---
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// --- Step 4: Fetch patientCustId ---
		var patientCustId string
		err := dbConn.Table(`"Users"`).
			Select(`"refUserCustId"`).
			Where(`"refUserId" = ?`, req.PatientId).
			Scan(&patientCustId).Error
		if err != nil || patientCustId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid patient ID or user not found.",
			})
			return
		}

		// --- Step 5: Fetch Scan Center Cust ID ---
		type ScanCenterResult struct {
			RefSCCustId string `gorm:"column:refSCCustId"`
		}
		var scanCenter ScanCenterResult
		err = dbConn.Table(`appointment."refAppointments" AS ra`).
			Joins(`JOIN public."ScanCenter" AS sc ON sc."refSCId" = ra."refSCId"`).
			Where(`ra."refAppointmentId" = ?`, req.AppointmentId).
			Scan(&scanCenter).Error
		if err != nil || scanCenter.RefSCCustId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid scan center configuration or appointment not found.",
			})
			return
		}

		// --- Step 6: Generate unique DICOM filename ---
		uniqueFilename := GenerateDicomFileName(
			scanCenter.RefSCCustId,
			patientCustId,
			req.Side,
			req.FileName,
		)
		s3Key := fmt.Sprintf("dicom/%s", uniqueFilename)

		// --- Step 7: Generate S3 Presigned URLs ---
		ctx := c.Request.Context()
		s3Client, err := s3config.New(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to initialize S3 client.",
			})
			return
		}

		uploadURL, err := s3Client.PresignPut(ctx, s3Key, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate presigned upload URL.",
			})
			return
		}

		viewURL, _ := s3Client.PresignGet(ctx, s3Key, 24*time.Hour)

		// --- Step 8: Prepare response payload ---
		payload := map[string]interface{}{
			"status":      true,
			"message":     "Presigned URLs generated successfully!",
			"uploadURL":   uploadURL,
			"viewURL":     viewURL,
			"s3Key":       s3Key,
			"fileName":    uniqueFilename,
			"oldFileName": req.FileName,
		}

		// --- Step 9: Generate token & respond ---
		token := accesstoken.CreateToken(idValue, roleIdValue)
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func SaveDicomController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		fmt.Println("Allowed")

		data, ok := helper.RequestHandler[model.SaveDicomReq](c)
		if !ok {
			return
		}

		// var reqVal model.AddIntakeFormReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.SaveDicomService(dbConn, *data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func ViewDicomController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		fmt.Println("Allowed")

		data, ok := helper.RequestHandler[model.ViewTechnicianIntakeFormReq](c)
		if !ok {
			return
		}

		// var reqVal model.AddIntakeFormReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		DicomData := service.ViewDicomService(dbConn, *data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":    true,
			"DicomData": DicomData,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func DeleteDicomController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.DeleteDicomReq](c, true)
		if !ok {
			return
		}

		// var reqVal model.AddIntakeFormReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.DeleteDicomService(dbConn, data)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func DownloadDicomFileController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID and role from context (set by auth middleware)
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID or RoleID not found in request context.",
			})
			return
		}

		// Parse and decrypt request body into DownloadDicomReq struct
		data, ok := helper.GetRequestBody[model.DownloadDicomReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var dicomFile model.DicomFileModel

		// Query file metadata by FileId
		err := dbConn.Raw(query.GetDicomFileSQL, data.FileId).Scan(&dicomFile).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch Dicom File: %v", err)
			payload := map[string]interface{}{
				"status":  false,
				"message": "Invalid Dicom File ID",
			}

			token := accesstoken.CreateToken(idValue, roleIdValue)

			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
			return // Important: stop further processing
		}

		filePath := "./Assets/Dicom/" + dicomFile.FileName
		fileName := dicomFile.FileName

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// Optionally, set Content-Length header
		if fi, err := os.Stat(filePath); err == nil {
			c.Header("Content-Length", fmt.Sprintf("%d", fi.Size()))
		}

		// Set headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		c.Header("Content-Type", "application/octet-stream")

		// Stream the file to the client
		c.File(filePath)
	}
}

// func DownloadMultipleDicomFilesController() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		data, ok := helper.GetRequestBody[model.OneDownloadDicomReq](c, true)
// 		if !ok {
// 			return
// 		}

// 		dbConn, sqlDB := db.InitDB()
// 		defer sqlDB.Close()

// 		var files []model.DicomFileModel
// 		err := dbConn.Raw(query.GetDicomFile, data.AppointmentId, data.UserId, data.Side).Scan(&files).Error
// 		if err != nil || len(files) == 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"status":  false,
// 				"message": "Failed to retrieve DICOM files.",
// 			})
// 			return
// 		}

// 		zipFilename := strings.Join(strings.Split(files[0].FileName, "_")[:len(strings.Split(files[0].FileName, "_"))-2], "_") + ".zip"
// 		// Set headers before writing data
// 		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilename)
// 		c.Writer.Header().Set("Content-Type", "application/zip")
// 		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
// 		c.Writer.Header().Set("Cache-Control", "no-cache")

// 		// Create zip writer directly on response writer
// 		zipWriter := zip.NewWriter(c.Writer)

// 		for _, file := range files {
// 			filePath := "./Assets/Dicom/" + file.FileName

// 			if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 				continue
// 			}

// 			fileToZip, err := os.Open(filePath)
// 			if err != nil {
// 				continue
// 			}

// 			writer, err := zipWriter.Create(file.FileName)
// 			if err != nil {
// 				fileToZip.Close()
// 				continue
// 			}

// 			_, err = io.Copy(writer, fileToZip)
// 			fileToZip.Close()
// 			if err != nil {
// 				continue
// 			}
// 		}

// 		// Close zip writer to flush all data to response
// 		if err := zipWriter.Close(); err != nil {
// 			log.Println("Error closing zip writer:", err)
// 		}

// 		// Flush the response writer to ensure all data is sent
// 		if flusher, ok := c.Writer.(http.Flusher); ok {
// 			flusher.Flush()
// 		}
// 	}
// }

func DownloadMultipleDicomFilesController() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, ok := helper.GetRequestBody[model.OneDownloadDicomReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var files []model.DicomFileModel
		err := dbConn.Raw(query.GetDicomFile, data.AppointmentId, data.UserId, data.Side).Scan(&files).Error
		if err != nil || len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Failed to retrieve DICOM files."})
			return
		}

		zipFilename := strings.Join(strings.Split(files[0].FileName, "_")[:len(strings.Split(files[0].FileName, "_"))-2], "_") + ".zip"

		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilename)
		c.Writer.Header().Set("Content-Type", "application/zip")
		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
		c.Writer.Header().Set("Cache-Control", "no-cache")

		zipWriter := zip.NewWriter(c.Writer)

		// ✅ Initialize AWS SDK config (only once)
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to load AWS config"})
			return
		}
		s3Client := s3.NewFromConfig(cfg)

		for _, file := range files {
			filePath := "./Assets/Dicom/" + file.FileName

			if strings.HasPrefix(file.FileName, "http") || strings.Contains(file.FileName, "amazonaws.com") {
				// ✅ Download from S3
				bucket, key := parseS3URL(file.FileName)
				obj, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
				})
				if err != nil {
					log.Println("❌ Failed to get object from S3:", err)
					continue
				}

				writer, err := zipWriter.Create(key[strings.LastIndex(key, "/")+1:])
				if err != nil {
					log.Println("❌ Failed to create zip entry:", err)
					obj.Body.Close()
					continue
				}

				_, err = io.Copy(writer, obj.Body)
				obj.Body.Close()
				if err != nil {
					log.Println("❌ Failed to copy S3 object:", err)
					continue
				}

			} else {
				// ✅ Local file handling (unchanged)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					continue
				}

				fileToZip, err := os.Open(filePath)
				if err != nil {
					continue
				}

				writer, err := zipWriter.Create(file.FileName)
				if err != nil {
					fileToZip.Close()
					continue
				}

				_, err = io.Copy(writer, fileToZip)
				fileToZip.Close()
				if err != nil {
					continue
				}
			}
		}

		if err := zipWriter.Close(); err != nil {
			log.Println("Error closing zip writer:", err)
		}

		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

// Helper to extract bucket and key from full S3 URL
func parseS3URL(url string) (bucket string, key string) {
	parts := strings.Split(url, ".s3.")
	bucket = strings.TrimPrefix(parts[0], "https://")
	keyParts := strings.SplitN(parts[1], "/", 2)
	key = keyParts[1]
	return
}

func AllDownloadDicomFileController() gin.HandlerFunc {
	return func(c *gin.Context) {

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		data, ok := helper.GetRequestBody[model.OneDownloadDicomReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var DicomFile model.DicomFileModel

		FileErr := dbConn.Raw(query.GetDicomFile, data.AppointmentId).Scan(&DicomFile).Error
		if FileErr != nil {
			log.Printf("ERROR: Failed to fetch Staff Available: %v", FileErr)
			payload := map[string]interface{}{
				"status":  false,
				"message": "Invalid Dicom File ID",
			}

			token := accesstoken.CreateToken(idValue, roleIdValue)

			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
		}

		filePath := "./Assets/Dicom/" + DicomFile.FileName
		fileName := DicomFile.FileName

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// Set headers
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/octet-stream")

		// Stream the file to response
		c.File(filePath)

	}
}

// func OverallDownloadDicomFileController() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		idValue, idExists := c.Get("id")
// 		roleIdValue, roleIdExists := c.Get("roleId")

// 		if !idExists || !roleIdExists {
// 			// Handle error: ID is missing from context (e.g., middleware didn't set it)
// 			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
// 				"status":  false,
// 				"message": "User ID, RoleID, Branch ID not found in request context.",
// 			})
// 			return // Stop processing
// 		}

// 		data, ok := helper.GetRequestBody[model.OverAllDicomModel](c, true)
// 		if !ok {
// 			return
// 		}

// 		dbConn, sqlDB := db.InitDB()
// 		defer sqlDB.Close()

// 		var files model.DicomFileModel

// 		FileErr := dbConn.Raw(query.GetAllDicomSQL, data.AppointmentId).Scan(&files).Error
// 		if FileErr != nil {
// 			log.Printf("ERROR: Failed to fetch Staff Available: %v", FileErr)
// 			payload := map[string]interface{}{
// 				"status":  false,
// 				"message": "Invalid Dicom File ID",
// 			}

// 			token := accesstoken.CreateToken(idValue, roleIdValue)

// 			c.JSON(http.StatusOK, gin.H{
// 				"data":  hashapi.Encrypt(payload, true, token),
// 				"token": token,
// 			})
// 		}

// 		zipFilename := "DicomFiles" + time.Now().Format("02-01-2006") + ".zip"
// 		// Set headers before writing data
// 		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilename)
// 		c.Writer.Header().Set("Content-Type", "application/zip")
// 		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
// 		c.Writer.Header().Set("Cache-Control", "no-cache")

// 		// Create zip writer directly on response writer
// 		zipWriter := zip.NewWriter(c.Writer)

// 		for _, file := range files.FileName {
// 		}

// 	}
// }

// func OverallDownloadDicomFileController() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		idValue, idExists := c.Get("id")
// 		roleIdValue, roleIdExists := c.Get("roleId")

// 		if !idExists || !roleIdExists {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"status":  false,
// 				"message": "User ID, RoleID, Branch ID not found in request context.",
// 			})
// 			return
// 		}

// 		data, ok := helper.GetRequestBody[model.OverAllDicomModel](c, true)
// 		if !ok {
// 			return
// 		}

// 		dbConn, sqlDB := db.InitDB()
// 		defer sqlDB.Close()

// 		// Fix: Use slice instead of single model
// 		var files []model.DicomFileModel

// 		FileErr := dbConn.Raw(query.GetAllDicomSQL, data.AppointmentId).Scan(&files).Error
// 		if FileErr != nil {
// 			log.Printf("ERROR: Failed to fetch DICOM Files: %v", FileErr)
// 			payload := map[string]interface{}{
// 				"status":  false,
// 				"message": "Invalid Dicom File ID",
// 			}

// 			token := accesstoken.CreateToken(idValue, roleIdValue)

// 			c.JSON(http.StatusOK, gin.H{
// 				"data":  hashapi.Encrypt(payload, true, token),
// 				"token": token,
// 			})
// 			return
// 		}

// 		// Check if files exist
// 		if len(files) == 0 {
// 			payload := map[string]interface{}{
// 				"status":  false,
// 				"message": "No DICOM files found",
// 			}

// 			token := accesstoken.CreateToken(idValue, roleIdValue)

// 			c.JSON(http.StatusOK, gin.H{
// 				"data":  hashapi.Encrypt(payload, true, token),
// 				"token": token,
// 			})
// 			return
// 		}

// 		// Create zip filename
// 		zipFilename := "DicomFiles_" + timeZone.GetTimeWithFormate("02-01-2006") + ".zip"

// 		// Set headers before writing data
// 		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilename)
// 		c.Writer.Header().Set("Content-Type", "application/zip")
// 		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
// 		c.Writer.Header().Set("Cache-Control", "no-cache")

// 		// Create zip writer directly on response writer
// 		zipWriter := zip.NewWriter(c.Writer)
// 		defer zipWriter.Close()

// 		// Process files by side (Right/Left)
// 		for _, file := range files {
// 			fmt.Println(file)

// 			// Determine folder name based on side
// 			var sideName string
// 			switch strings.ToLower(file.Side) {
// 			case "right", "r":
// 				sideName = "Right"
// 			case "left", "l":
// 				sideName = "Left"
// 			default:
// 				sideName = "Other" // For files without specific side
// 			}

// 			// Extract the base pattern from filename and replace L/R with full side name
// 			fileName := file.FileName

// 			// Remove the file extension
// 			nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

// 			// Split by underscore and remove the last part (sequence number)
// 			parts := strings.Split(nameWithoutExt, "_")
// 			if len(parts) > 1 {
// 				// Remove the last part (sequence number like "1", "2", etc.)
// 				parts = parts[:len(parts)-1]
// 			}

// 			// Join back to create base pattern
// 			basePattern := strings.Join(parts, "_")

// 			// Replace L/R with full side name in the base pattern
// 			if strings.HasSuffix(basePattern, "_L") {
// 				basePattern = strings.TrimSuffix(basePattern, "_L")
// 			} else if strings.HasSuffix(basePattern, "_R") {
// 				basePattern = strings.TrimSuffix(basePattern, "_R")
// 			}

// 			// Create folder name: "BasePattern_Side"
// 			folderName := fmt.Sprintf("%s_%s", basePattern, sideName)

// 			// Create file path in zip: "BasePattern_Side/FileName"
// 			zipPath := fmt.Sprintf("%s/%s", folderName, file.FileName)

// 			// Read the actual file from disk
// 			filePath := filepath.Join("./Assets/Dicom/", file.FileName)
// 			fileData, err := os.ReadFile(filePath)
// 			if err != nil {
// 				log.Printf("ERROR: Failed to read file %s: %v", file.FileName, err)
// 				continue // Skip this file and continue with others
// 			}

// 			// Create file in zip
// 			zipFile, err := zipWriter.Create(zipPath)
// 			if err != nil {
// 				log.Printf("ERROR: Failed to create zip entry for %s: %v", file.FileName, err)
// 				continue
// 			}

// 			// Write file content to zip
// 			_, err = zipFile.Write(fileData)
// 			if err != nil {
// 				log.Printf("ERROR: Failed to write file %s to zip: %v", file.FileName, err)
// 				continue
// 			}

// 			log.Printf("Added file %s to %s folder in zip", file.FileName, folderName)
// 		}

// 		// Zip writer will be closed by defer, and response will be sent
// 		log.Printf("Zip file %s created successfully with %d files", zipFilename, len(files))
// 	}
// }

// func OverallDownloadDicomFileController() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		idValue, idExists := c.Get("id")
// 		roleIdValue, roleIdExists := c.Get("roleId")

// 		if !idExists || !roleIdExists {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"status":  false,
// 				"message": "User ID, RoleID, Branch ID not found in request context.",
// 			})
// 			return
// 		}

// 		data, ok := helper.GetRequestBody[model.OverAllDicomModel](c, true)
// 		if !ok {
// 			return
// 		}

// 		dbConn, sqlDB := db.InitDB()
// 		defer sqlDB.Close()

// 		var files []model.DicomFileModel
// 		FileErr := dbConn.Raw(query.GetAllDicomSQL, data.AppointmentId).Scan(&files).Error
// 		if FileErr != nil {
// 			log.Printf("ERROR: Failed to fetch DICOM Files: %v", FileErr)
// 			payload := map[string]interface{}{
// 				"status":  false,
// 				"message": "Invalid Dicom File ID",
// 			}
// 			token := accesstoken.CreateToken(idValue, roleIdValue)
// 			c.JSON(http.StatusOK, gin.H{
// 				"data":  hashapi.Encrypt(payload, true, token),
// 				"token": token,
// 			})
// 			return
// 		}

// 		if len(files) == 0 {
// 			payload := map[string]interface{}{
// 				"status":  false,
// 				"message": "No DICOM files found",
// 			}
// 			token := accesstoken.CreateToken(idValue, roleIdValue)
// 			c.JSON(http.StatusOK, gin.H{
// 				"data":  hashapi.Encrypt(payload, true, token),
// 				"token": token,
// 			})
// 			return
// 		}

// 		zipFilename := "DicomFiles_" + timeZone.GetTimeWithFormate("02-01-2006") + ".zip"
// 		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilename)
// 		c.Writer.Header().Set("Content-Type", "application/zip")
// 		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
// 		c.Writer.Header().Set("Cache-Control", "no-cache")

// 		zipWriter := zip.NewWriter(c.Writer)
// 		defer zipWriter.Close()

// 		// S3 client (if needed)
// 		ctx := c.Request.Context()
// 		s3Client, _ := s3config.New(ctx) // handle error if you want

// 		for _, file := range files {
// 			var fileData []byte
// 			var err error

// 			if strings.HasPrefix(file.FileName, "http") {
// 				// Fetch file from S3
// 				fileData, err = s3Client.DownloadBytes(ctx, file.FileName) // Implement DownloadBytes() in your s3 helper
// 				if err != nil {
// 					log.Printf("ERROR: Failed to download S3 file %s: %v", file.FileName, err)
// 					continue
// 				}
// 			} else {
// 				// Fetch file from local disk
// 				filePath := filepath.Join("./Assets/Dicom/", file.FileName)
// 				fileData, err = os.ReadFile(filePath)
// 				if err != nil {
// 					log.Printf("ERROR: Failed to read local file %s: %v", file.FileName, err)
// 					continue
// 				}
// 			}

// 			// Determine folder by side
// 			sideName := "Other"
// 			switch strings.ToLower(file.Side) {
// 			case "right", "r":
// 				sideName = "Right"
// 			case "left", "l":
// 				sideName = "Left"
// 			}

// 			nameWithoutExt := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
// 			parts := strings.Split(nameWithoutExt, "_")
// 			if len(parts) > 1 {
// 				parts = parts[:len(parts)-1]
// 			}
// 			basePattern := strings.Join(parts, "_")
// 			if strings.HasSuffix(basePattern, "_L") {
// 				basePattern = strings.TrimSuffix(basePattern, "_L")
// 			} else if strings.HasSuffix(basePattern, "_R") {
// 				basePattern = strings.TrimSuffix(basePattern, "_R")
// 			}

// 			folderName := fmt.Sprintf("%s_%s", basePattern, sideName)
// 			zipPath := fmt.Sprintf("%s/%s", folderName, filepath.Base(file.FileName))

// 			zipFile, err := zipWriter.Create(zipPath)
// 			if err != nil {
// 				log.Printf("ERROR: Failed to create zip entry for %s: %v", file.FileName, err)
// 				continue
// 			}

// 			_, err = zipFile.Write(fileData)
// 			if err != nil {
// 				log.Printf("ERROR: Failed to write file %s to zip: %v", file.FileName, err)
// 				continue
// 			}

// 			log.Printf("Added file %s to %s folder in zip", file.FileName, folderName)
// 		}

// 		log.Printf("Zip file %s created successfully with %d files", zipFilename, len(files))
// 	}
// }

func OverallDownloadDicomFileController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID or Role ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.OverAllDicomModel](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var files []model.DicomFileModel
		err := dbConn.Raw(query.GetAllDicomSQL, data.AppointmentId).Scan(&files).Error
		if err != nil {
			log.Printf("ERROR: Failed to fetch DICOM Files: %v", err)
			payload := map[string]interface{}{
				"status":  false,
				"message": "Invalid Dicom File ID",
			}
			token := accesstoken.CreateToken(idValue, roleIdValue)
			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
			return
		}

		if len(files) == 0 {
			payload := map[string]interface{}{
				"status":  false,
				"message": "No DICOM files found",
			}
			token := accesstoken.CreateToken(idValue, roleIdValue)
			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
			return
		}

		// ✅ Initialize S3 client
		ctx := c.Request.Context()

		groupedFiles := make(map[string][]map[string]string)

		for _, file := range files {
			// Determine side folder
			sideName := "Other"
			switch strings.ToLower(file.Side) {
			case "right", "r":
				sideName = "Right"
			case "left", "l":
				sideName = "Left"
			}

			// Generate base pattern from filename
			nameWithoutExt := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
			parts := strings.Split(nameWithoutExt, "_")
			if len(parts) > 1 {
				parts = parts[:len(parts)-1]
			}
			basePattern := strings.Join(parts, "_")

			if strings.HasSuffix(basePattern, "_L") {
				basePattern = strings.TrimSuffix(basePattern, "_L")
			} else if strings.HasSuffix(basePattern, "_R") {
				basePattern = strings.TrimSuffix(basePattern, "_R")
			}

			folderName := fmt.Sprintf("%s_%s", basePattern, sideName)

			// ✅ Generate URL (either S3 presigned or local)
			fileURL := file.FileName
			if strings.HasPrefix(file.FileName, "http") {
				// Generate pre-signed S3 URL (1-hour validity)
				presignedURL, err := s3Service.GeneratePresignGetURL(ctx, file.FileName, time.Hour)
				if err != nil {
					log.Printf("ERROR: Failed to generate presigned URL for %s: %v", file.FileName, err)
					continue
				}
				fileURL = presignedURL
			} else {
				// Local path (BASE_URL + /Assets/Dicom/)
				fileURL = fmt.Sprintf("%s/Assets/Dicom/%s", os.Getenv("BASE_URL"), file.FileName)
			}

			groupedFiles[folderName] = append(groupedFiles[folderName], map[string]string{
				"fileName": file.FileName,
				"url":      fileURL,
				"side":     file.Side,
			})
		}

		payload := map[string]interface{}{
			"status":  true,
			"message": fmt.Sprintf("%d DICOM files found", len(files)),
			"folders": groupedFiles, // 👈 grouped by folder name
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
