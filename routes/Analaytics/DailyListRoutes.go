package routes

import (
	controllers "AuthenticationService/controllers/Analaytics"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitDailyListRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/dailyList")
	route.POST("/", accesstoken.JWTMiddleware(), controllers.GetDailyListController())
}
