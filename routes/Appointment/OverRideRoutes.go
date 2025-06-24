package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitOverrideRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/override")
	route.POST("/list", accesstoken.JWTMiddleware(), controllers.ListOverrideController())
	route.POST("/update", accesstoken.JWTMiddleware(), controllers.WriteOverrideController())
}
