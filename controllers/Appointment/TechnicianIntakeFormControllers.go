package controllers

import (
	service "AuthenticationService/Service/Appointment"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Appointment"
	query "AuthenticationService/query/Appointment"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

		ViewIntakeData, Aduit, TechIntakeData := service.ViewTechnicianIntakeFormService(dbConn, *data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"IntakeData":      ViewIntakeData,
			"IntakeDataAduit": Aduit,
			"TechIntakeData":  TechIntakeData,
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
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		uploadPath := "./Assets/Dicom/"

		log := logger.InitLogger()

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Error retrieving profile image from request: " + err.Error(),
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
		// if ext != ".zip" {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"status":  false,
		// 		"message": "Invalid file type. Only .zip files are allowed.",
		// 	})
		// 	return
		// }

		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),                 // Generate a random UUID
			time.Now().Format("20060102150405"), // Add timestamp (YYYYMMDDHHMMSS)
			ext)                                 // Keep original file extension
		destinationPath := filepath.Join(uploadPath, uniqueFilename)

		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			log.Printf("Error creating upload directory '%s': %v\n", uploadPath, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Server error: Could not prepare image storage.",
			})
			return
		}

		if err := c.SaveUploadedFile(file, destinationPath); err != nil {
			log.Printf("Error saving uploaded file to '%s': %v\n", destinationPath, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Server error: Could not save profile image.",
			})
			return
		}

		log.Printf("Successfully uploaded image: %s\n", destinationPath)

		payload := map[string]interface{}{
			"status":      true,
			"message":     "Dicom File uploaded successfully!",
			"fileName":    uniqueFilename,
			"oldFilename": file.Filename,
		}

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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed to retrieve DICOM files.",
			})
			return
		}

		zipFilename := strings.Join(strings.Split(files[0].FileName, "_")[:len(strings.Split(files[0].FileName, "_"))-2], "_") + ".zip"
		// Set headers before writing data
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipFilename)
		c.Writer.Header().Set("Content-Type", "application/zip")
		c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
		c.Writer.Header().Set("Cache-Control", "no-cache")

		// Create zip writer directly on response writer
		zipWriter := zip.NewWriter(c.Writer)

		for _, file := range files {
			filePath := "./Assets/Dicom/" + file.FileName

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

		// Close zip writer to flush all data to response
		if err := zipWriter.Close(); err != nil {
			log.Println("Error closing zip writer:", err)
		}

		// Flush the response writer to ensure all data is sent
		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}
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
