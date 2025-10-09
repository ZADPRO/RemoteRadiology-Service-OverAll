package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitRadiologistRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/radiologist")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostRadiologistController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchRadiologistController())
}
