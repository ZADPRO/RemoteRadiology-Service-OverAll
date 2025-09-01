package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitIntakeFormRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/intakeform")
	route.POST("/add", accesstoken.JWTMiddleware(), controllers.AddIntakeFormController())
	route.POST("/view", accesstoken.JWTMiddleware(), controllers.ViewIntakeFormController())
	route.POST("/verify", accesstoken.JWTMiddleware(), controllers.VerifyIntakeFormController())
	route.POST("/update", accesstoken.JWTMiddleware(), controllers.UpdateIntakeFormController())
	route.POST("/getReportData", accesstoken.JWTMiddleware(), controllers.GetReportDataController())
	route.POST("/getConsentData", accesstoken.JWTMiddleware(), controllers.GetConsentDataController())
	route.POST("/allowoverride", accesstoken.JWTMiddleware(), controllers.AllowOverrideController())
}
