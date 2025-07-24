package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitWellthgreenPerformingProviderRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/wellgreenperformingprovider")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostWellgreenPerformingProviderController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchWellgreenPerformingProviderController())
}
