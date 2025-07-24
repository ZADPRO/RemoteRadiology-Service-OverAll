package controllers

import (
	service "AuthenticationService/Service/Analaytics"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	hashdb "AuthenticationService/internal/Helper/HashDB"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Analaytics"
	query "AuthenticationService/query/Analaytics"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func AddTrainingMaterialController() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.AddTrainingMaterialReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AddTrainingMaterialService(dbConn, data, int(idValue.(float64)))

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

func ListTrainingMaterialController() gin.HandlerFunc {
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

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		Data := service.ListTrainingMaterialService(dbConn, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status": true,
			"data":   Data,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func DeleteTrainingMaterialController() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.DeleteTrainingMaterialReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.DeleteTrainingMaterialService(dbConn, data, int(idValue.(float64)))

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

func DownloadTrainingMaterialController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID or RoleID not found in context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.DeleteTrainingMaterialReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var FileData []model.CreateMaterialModel
		err := dbConn.Raw(query.OneListTrainingFilesSQL, data.Id).Scan(&FileData).Error
		if err != nil || len(FileData) == 0 {
			log.Printf("ERROR: Failed to fetch file metadata: %v", err)
			payload := map[string]interface{}{
				"status":  false,
				"message": "Invalid or missing Training Material",
			}
			token := accesstoken.CreateToken(idValue, roleIdValue)
			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
			return
		}

		// Decrypt fields
		for i, f := range FileData {
			FileData[i].TMFileName = hashdb.Decrypt(f.TMFileName)
			FileData[i].TMFilePath = hashdb.Decrypt(f.TMFilePath)
			fmt.Println(FileData[i])
		}

		fmt.Println("&&&&&&&&&&&&&&&&&&&&", FileData[0].TMFilePath)

		filePath := "./Assets/Files/" + FileData[0].TMFilePath
		fileName := FileData[0].TMFileName

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// Set headers
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")

		// This will show the PDF in an <iframe> or browser tab:
		c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, url.QueryEscape(fileName)))
		c.Header("Content-Type", "application/pdf")

		// Stream file
		c.File(filePath)
	}
}
