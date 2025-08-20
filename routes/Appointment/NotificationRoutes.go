package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitNotificationRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/notification")
	route.POST("/viewnotification", accesstoken.JWTMiddleware(), controllers.NotificationController())
	route.POST("/readStatus", accesstoken.JWTMiddleware(), controllers.ReadStatusController())
	route.GET("/getUnreadCount", accesstoken.JWTMiddleware(), controllers.GetNotificationCountController())
}
