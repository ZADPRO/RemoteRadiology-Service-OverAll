package controllers

import (
	service "AuthenticationService/Service/Appointment"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Appointment"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListOverrideController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.OverrideListReq](c)
		if !ok {
			return
		}

		// var reqVal model.OverrideListReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		ListAllOverrideData := service.ListAllOverRideService(dbConn, *data)

		payload := map[string]interface{}{
			"data": ListAllOverrideData,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func WriteOverrideController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.WriteOverrideListReq](c)
		if !ok {
			return
		}

		// var reqVal model.OverrideListReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.WriteOverRideService(dbConn, *data, int(idValue.(float64)))

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
