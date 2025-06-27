package controllers

import (
	service "AuthenticationService/Service/Appointment"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Appointment"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
		if ext != ".zip" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid file type. Only .zip files are allowed.",
			})
			return
		}

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
