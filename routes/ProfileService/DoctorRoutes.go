package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitDoctorRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/doctor")
	route.POST("/list-alldoctor", accesstoken.JWTMiddleware(), controllers.GetAllDoctorDataController())
	route.POST("/list-doctor", accesstoken.JWTMiddleware(), controllers.GetDoctorDataController())
}
