package controllers

import (
	service "AuthenticationService/Service/Analaytics"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Analaytics"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOverallAnalayticsController() gin.HandlerFunc {
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

		Value := service.AdminOverallAnalayticsService(dbConn)

		payload := map[string]interface{}{
			"status":                              true,
			"AdminOverallAnalaytics":              Value.AdminScanCenterModel,
			"AdminOverallScanIndicatesAnalaytics": Value.AdminOverallScanIndicatesAnalayticsModel,
			"AllScaCenter":                        Value.GetAllScaCenter,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AdminOverallOneAnalayticsController() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.AdminOverallOneAnalyticsReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		Value := service.AdminOverallOneAnalayticsService(dbConn, data)

		payload := map[string]interface{}{
			"status":                              true,
			"AdminOverallAnalaytics":              Value.AdminScanCenterModel,
			"AdminOverallScanIndicatesAnalaytics": Value.AdminOverallScanIndicatesAnalayticsModel,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
