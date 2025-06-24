package controllers

import (
	service "AuthenticationService/Service/UserService"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/UserService"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostDoctorController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.DoctorRegisterReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PostDoctorService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func PatchDoctorController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.UpdateDoctorReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PatchDoctorService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
