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

func AdminOverallOneAnalayticsController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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
		data, ok := helper.GetRequestBody[model.AdminOverallOneAnalyticsReq](c, false)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		Value := service.AdminOverallOneAnalayticsService(dbConn, data, int(roleIdValue.(float64)))

		//Load the Payload
		payload := map[string]interface{}{
			"status":                              true,
			"AdminOverallAnalaytics":              Value.AdminScanCenterModel,
			"AdminOverallScanIndicatesAnalaytics": Value.AdminOverallScanIndicatesAnalayticsModel,
			"AllScaCenter":                        Value.GetAllScaCenter,
			"UserListIds":                         Value.UserListIdsData,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, false, token),
			"token": token,
		})
	}
}

func OneUserController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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
		data, ok := helper.GetRequestBody[model.OneUserReq](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		Value := service.OneUserService(dbConn, data, int(idValue.(float64)), int(roleIdValue.(float64)))

		//Load the Payload
		payload := map[string]interface{}{
			"status":                              true,
			"AdminOverallAnalaytics":              Value.AdminScanCenterModel,
			"AdminOverallScanIndicatesAnalaytics": Value.AdminOverallScanIndicatesAnalayticsModel,
			"UserAccessTiming":                    Value.UserAccessTimingModel,
			"ListScanAppointmentCount":            Value.ListScanAppointmentCountModel,
			"TotalCorrectEdit":                    Value.TotalCorrectEdit,
			"ImpressionModel":                     Value.ImpressionModel,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
