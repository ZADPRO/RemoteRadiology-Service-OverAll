package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitRegisterPatientRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/patient")
	route.POST("/check", accesstoken.JWTMiddleware(), controllers.PostCheckPatientController())
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostAddPatientController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchUpdatePatientController())
	route.PATCH("/createAppointment", accesstoken.JWTMiddleware(), controllers.PostCreatePatientController())
	route.PATCH("/sendMailAppointment", accesstoken.JWTMiddleware(), controllers.PostSendMailPatientController())
}
