package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitManageAppointmentRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/manageappointment")
	route.POST("/add", accesstoken.JWTMiddleware(), controllers.AddAppointmentController())
	route.GET("/viewpatienthistory", accesstoken.JWTMiddleware(), controllers.ViewPatientHistory())
	route.POST("/viewtechnicianpatientqueue", accesstoken.JWTMiddleware(), controllers.ViewTechnicianPatientQueue())
}
