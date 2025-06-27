package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitTechnicianIntakeFormRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/technicianintakeform")
	route.POST("/add", accesstoken.JWTMiddleware(), controllers.AddTechnicianIntakeFormController())
	route.POST("/dicomupload", accesstoken.JWTMiddleware(), controllers.PostUploadDicomFileController())
}
