package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitTechnicianIntakeFormRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/technicianintakeform")
	route.POST("/add", accesstoken.JWTMiddleware(), controllers.AddTechnicianIntakeFormController())
	route.POST("/view", accesstoken.JWTMiddleware(), controllers.ViewTechnicianIntakeFormController())
	route.POST("/dicomupload", accesstoken.JWTMiddleware(), controllers.PostUploadDicomFileController())
	route.POST("/savedicom", accesstoken.JWTMiddleware(), controllers.SaveDicomController())
	route.POST("/viewDicom", accesstoken.JWTMiddleware(), controllers.ViewDicomController())
	route.POST("/downloaddicom", accesstoken.JWTMiddleware(), controllers.DownloadDicomFileController())
	route.POST("/alldownloaddicom", accesstoken.JWTMiddleware(), controllers.DownloadMultipleDicomFilesController())
}
