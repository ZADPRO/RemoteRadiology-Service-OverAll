package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"

)

func InitImageRoutes(router *gin.Engine) {
	route := router.Group("/api/v1")
	route.POST("/upload-profile-image", accesstoken.JWTMiddleware(), controllers.PostUploadProfileImage())
	route.POST("/upload-public-profile-image", accesstoken.JWTMiddleware(), controllers.PostUploadPublicProfileImage())
	route.POST("/upload-private-document", accesstoken.JWTMiddleware(), controllers.PostUploadPrivateDocument())


}
