package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitTechnicianRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/technician")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostAddTechnicianController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchUpdateTechnicianController())
	// route.POST("/map", accesstoken.JWTMiddleware(), controllers.PostManageTechnicianMapController())
	// route.GET("/", accesstoken.JWTMiddleware(), controllers.GetTechnicianController())
}
