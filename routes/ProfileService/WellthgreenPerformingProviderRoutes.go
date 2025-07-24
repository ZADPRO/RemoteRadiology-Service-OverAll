package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitWellthgreenPerformingProviderRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/wellthgreenperformingprovider")
	route.GET("/list-allperformingprovider", accesstoken.JWTMiddleware(), controllers.GetAllPerformingProviderDataController())
	route.POST("/list-performingprovider", accesstoken.JWTMiddleware(), controllers.GetPerformingProviderDataController())
}
