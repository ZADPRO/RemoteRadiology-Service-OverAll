package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitReceptionistRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/receptionist")
	route.POST("/new", accesstoken.JWTMiddleware(), controllers.PostAddReceptionistController())
	route.PATCH("/update", accesstoken.JWTMiddleware(), controllers.PatchReceptionistController())
	// route.POST("/map", accesstoken.JWTMiddleware(), controllers.PostManageReceptionistMapController())
	// route.GET("/", accesstoken.JWTMiddleware(), controllers.GetTechnicianController())
}
