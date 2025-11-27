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

func GetDailyListController() gin.HandlerFunc {
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
		data, ok := helper.GetRequestBody[model.GetDailyListModel](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		Status, resVal := service.GetDailyListService(dbConn, data)

		//Load the Payload
		payload := map[string]interface{}{
			"status": Status,
			"data":   resVal,
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
