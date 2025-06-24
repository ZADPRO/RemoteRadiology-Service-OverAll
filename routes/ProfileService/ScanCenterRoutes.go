package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitScanCenterRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/scancenter")
	route.GET("/list-allscan-center", accesstoken.JWTMiddleware(), controllers.GetAllScanCenterController())
	route.POST("/list-scan-center", accesstoken.JWTMiddleware(), controllers.GetScanCenterController())
}
