package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitPatientRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/patient")
	route.POST("/list-allpatient", accesstoken.JWTMiddleware(), controllers.GetAllPatientController())
	route.POST("/list-patient", accesstoken.JWTMiddleware(), controllers.GetPatientDataController())
}
