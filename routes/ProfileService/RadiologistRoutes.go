package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitRadiologistRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/radiologist")
	route.GET("/list-allradiologist", accesstoken.JWTMiddleware(), controllers.GetAllRadiologistDataController())
	route.POST("/list-radiologist", accesstoken.JWTMiddleware(), controllers.GetRadiologistDataController())
}
