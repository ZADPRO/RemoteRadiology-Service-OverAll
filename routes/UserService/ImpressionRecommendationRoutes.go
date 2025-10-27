package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitImpressionRecommendationRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/impressionrecommendation")
	route.GET("/", accesstoken.JWTMiddleware(), controllers.GetImpressionRecommendationController())
	route.POST("/add", accesstoken.JWTMiddleware(), controllers.AddImpressionRecommendationController())
	route.POST("/update", accesstoken.JWTMiddleware(), controllers.UpdateImpressionRecommendationController())
	route.POST("/delete", accesstoken.JWTMiddleware(), controllers.DeleteImpressionRecommendationController())
	route.POST("/updateorder", accesstoken.JWTMiddleware(), controllers.UpdateOrderImpressionRecommendationController())
}
