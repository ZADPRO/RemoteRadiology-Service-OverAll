package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitManagerRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/manager")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostManagerController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchManagerController())
}
