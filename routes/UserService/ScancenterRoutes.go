package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitScanCenterRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/scancenter")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostScanCenterController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchScanCenterController())
}
