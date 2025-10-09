package controllers

import (
	service "AuthenticationService/Service/Authentication"
	db "AuthenticationService/internal/DB"
	model "AuthenticationService/internal/Model/Authentication"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForgetPasswordController() gin.HandlerFunc {

	return func(c *gin.Context) {

		var reqVal model.ForgetPasswordReq

		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try Again " + err.Error(),
			})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ForgetPasswordService(dbConn, reqVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		// if resVal.Status {
		// 	response["token"] = resVal.Token
		// 	response["RoleType"] = resVal.RoleType
		// }

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})

	}
}

func VerifyForgetPasswordOTPController() gin.HandlerFunc {

	return func(c *gin.Context) {

		var reqVal model.VerifyPasswordReq

		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try Again " + err.Error(),
			})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.VerifyForgetPasswordOTPService(dbConn, reqVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		// if resVal.Status {
		// 	response["token"] = resVal.Token
		// 	response["RoleType"] = resVal.RoleType
		// }

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})

	}
}

func ChangePasswordController() gin.HandlerFunc {

	return func(c *gin.Context) {

		var reqVal model.ChangePasswordReq

		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try Again " + err.Error(),
			})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ChangePasswordService(dbConn, reqVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		// if resVal.Status {
		// 	response["token"] = resVal.Token
		// 	response["RoleType"] = resVal.RoleType
		// }

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})

	}

}
