package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitCoDoctorRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/codoctor")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostCoDoctorController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchCoDoctorController())
	// route.GET("/", accesstoken.JWTMiddleware(), controllers.GetTechnicianController())
}
