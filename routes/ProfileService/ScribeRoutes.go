package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitScribeRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/doctor")
	route.GET("/list-allScribe", accesstoken.JWTMiddleware(), controllers.GetAllScribeDataController())
	route.POST("/list-scribe", accesstoken.JWTMiddleware(), controllers.GetScribeDataController())
}
