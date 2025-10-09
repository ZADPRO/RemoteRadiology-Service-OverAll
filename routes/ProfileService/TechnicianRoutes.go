package routes

import (
	controllers "AuthenticationService/controllers/ProfileService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitTechnicianRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/profile/technician")
	route.POST("/list-alltechnician", accesstoken.JWTMiddleware(), controllers.GetAllTechnicianDataController())
	route.POST("/list-technician", accesstoken.JWTMiddleware(), controllers.GetOneTechnicianDataController())
}
