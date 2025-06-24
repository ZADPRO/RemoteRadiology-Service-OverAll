package routes

import (
	controllers "AuthenticationService/controllers/Authentication"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitLoginRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/authentication")
	route.POST("/login", controllers.LoginController())
	route.POST("/verifyotp", controllers.VerifyOTPController())
	route.POST("/changepassword", accesstoken.JWTMiddleware(), controllers.UserChangePasswordController())
}
