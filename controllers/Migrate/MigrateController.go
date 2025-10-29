package controllers

import (
	service "AuthenticationService/Service/Migrate"
	db "AuthenticationService/internal/DB"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DicomMigrateController() gin.HandlerFunc {

	return func(c *gin.Context) {

		// var reqVal model.ForgetPasswordReq

		// if err := c.BindJSON(&reqVal); err != nil {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"status":  false,
		// 		"message": "Something went wrong, Try Again " + err.Error(),
		// 	})
		// 	return
		// }

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.DicomMigrateService(dbConn)

		response := gin.H{
			"status":  status,
			"message": message,
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
