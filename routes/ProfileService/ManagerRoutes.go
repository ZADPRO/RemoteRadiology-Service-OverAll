package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitManagerRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/manager")
	route.GET("/list-allmanager", accesstoken.JWTMiddleware(), controllers.GetAllManagerDataController())
	route.POST("/list-manager", accesstoken.JWTMiddleware(), controllers.GetManagerDataController())
}
