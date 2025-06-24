package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitCoDoctorRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/codoctor")
	route.POST("/list-allcodoctor", accesstoken.JWTMiddleware(), controllers.GetAllCoDoctorDataController())
	route.POST("/list-codoctor", accesstoken.JWTMiddleware(), controllers.GetCoDoctorDataController())
}
