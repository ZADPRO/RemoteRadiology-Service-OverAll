package routes

import (
	controllers "AuthenticationService/controllers/UserService"

	"github.com/gin-gonic/gin"
)

func InitPatientRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/patient")
	route.POST("/getOtp", controllers.PostGetOtpPatientController())
	route.POST("/verifyOtp", controllers.PostCheckOTPPatientController())
	route.POST("/register", controllers.PostRegisterPatientController())
	// route.GET("/", accesstoken.JWTMiddleware(), controllers.GetTechnicianController())
}
