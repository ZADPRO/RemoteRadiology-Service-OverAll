package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitReceptionistRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/receptionist")
	route.POST("/list-allreceptionists", accesstoken.JWTMiddleware(), controllers.GetAllReceptionistDataController())
	route.POST("/list-receptionists", accesstoken.JWTMiddleware(), controllers.GetOneReceptionistDataController())
}
