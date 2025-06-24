package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitFilesRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/")
	route.POST("/upload-file", accesstoken.JWTMiddleware(), controllers.PostUploadFileController())
}
