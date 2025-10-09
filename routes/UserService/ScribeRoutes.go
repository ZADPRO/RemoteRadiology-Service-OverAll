package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitScribeRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/scribe")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostScribeController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchScribeController())
}
