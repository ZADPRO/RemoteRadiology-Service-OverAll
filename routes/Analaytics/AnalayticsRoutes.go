package routes

import (
	controllers "AuthenticationService/controllers/Analaytics"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitAnalayticsRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/analaytics")
	route.POST("/admin/overallscancenter", accesstoken.JWTMiddleware(), controllers.AdminOverallAnalayticsController())
	route.POST("/admin/overallonescancenter", accesstoken.JWTMiddleware(), controllers.AdminOverallOneAnalayticsController())
}
