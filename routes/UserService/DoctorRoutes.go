package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitDoctorRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/doctor")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostDoctorController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchDoctorController())
}
