package controllers

import (
	service "AuthenticationService/Service/UserService"
	db "AuthenticationService/internal/DB"
	model "AuthenticationService/internal/Model/UserService"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostGetOtpPatientController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.GetOtpPatient

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request body: " + err.Error(),
			})
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.PostGetOtpPatientService(dbConn, data)

		c.JSON(http.StatusOK, gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		})

	}
}

func PostCheckOTPPatientController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.VerifyOtpPatient

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request body: " + err.Error(),
			})
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.PostCheckOTPPatientService(dbConn, data)

		c.JSON(http.StatusOK, gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		})

	}
}

func PostRegisterPatientController() gin.HandlerFunc {
	return func(c *gin.Context) {

		var data model.RegisterPatientReq

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request body: " + err.Error(),
			})
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.PostRegisterPatientService(dbConn, data)

		// payload := map[string]interface{}{
		// 	"status":  resVal.Status,
		// 	"message": resVal.Message,
		// }

		c.JSON(http.StatusOK, gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		})

	}
}
