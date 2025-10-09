package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/user")
	route.GET("/", accesstoken.JWTMiddleware(), controllers.GetUserController())
	route.GET("/dashboard", accesstoken.JWTMiddleware(), controllers.DashboardController())
}
