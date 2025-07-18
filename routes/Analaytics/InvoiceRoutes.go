package routes

import (
	controllers "AuthenticationService/controllers/Analaytics"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitInvoiceRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/invoice")
	route.GET("/getamount", accesstoken.JWTMiddleware(), controllers.GetAmountController())
	route.POST("/updateamount", accesstoken.JWTMiddleware(), controllers.UpdateAmountController())
	route.POST("/getInvoiceData", accesstoken.JWTMiddleware(), controllers.GetInvoiceDataController())
	route.POST("/generteInvoice", accesstoken.JWTMiddleware(), controllers.GenerateInvoiceDataController())
	route.POST("/getInvoiceHistory", accesstoken.JWTMiddleware(), controllers.GetInvoiceHistoryController())
}
