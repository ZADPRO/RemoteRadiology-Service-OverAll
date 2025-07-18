package routes

import (
	controllers "AuthenticationService/controllers/Analaytics"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitTrainingMaterialRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/trainingmaterial")
	route.POST("/add", accesstoken.JWTMiddleware(), controllers.AddTrainingMaterialController())
	route.GET("/list", accesstoken.JWTMiddleware(), controllers.ListTrainingMaterialController())
	route.POST("/delete", accesstoken.JWTMiddleware(), controllers.DeleteTrainingMaterialController())
	route.POST("/download", accesstoken.JWTMiddleware(), controllers.DownloadTrainingMaterialController())
}
