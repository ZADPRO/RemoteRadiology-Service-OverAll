package controllers

import (
	service "AuthenticationService/Service/UserService"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/UserService"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetImpressionRecommendationController() gin.HandlerFunc {
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

		// data, ok := helper.RequestHandler[model.DoctorRegisterReq](c)
		// if !ok {
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		categoryData, ImpressionRecommendation := service.GetImpressionRecommendationService(dbConn)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":                   true,
			"categoryData":             categoryData,
			"ImpressionRecommendation": ImpressionRecommendation,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func AddImpressionRecommendationController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.AddImpressionRecommendationReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AddImpressionRecommendationService(dbConn, *data)
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

func UpdateImpressionRecommendationController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.UpdateImpressionRecommendationReq](c)
		if !ok {
			return
		}

		fmt.Println("=========>", data)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.UpdateImpressionRecommendationService(dbConn, *data)
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

func DeleteImpressionRecommendationController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.UpdateImpressionRecommendationReq](c)
		if !ok {
			return
		}

		fmt.Println("=========>", data)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.DeleteImpressionRecommendationService(dbConn, *data)
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

func UpdateOrderImpressionRecommendationController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.UpdateOrderImpressionRecommendationReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.UpdateOrderImpressionRecommendationService(dbConn, *data)
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
