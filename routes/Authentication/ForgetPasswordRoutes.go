package routes

import (
	controllers "AuthenticationService/controllers/Authentication"

	"github.com/gin-gonic/gin"
)

func InitForgetPasswordRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/forgetpassword")
	route.POST("/verifyuser", controllers.ForgetPasswordController())
	route.POST("/verifyotp", controllers.VerifyForgetPasswordOTPController())
	route.POST("/changepassword", controllers.ChangePasswordController())
}
