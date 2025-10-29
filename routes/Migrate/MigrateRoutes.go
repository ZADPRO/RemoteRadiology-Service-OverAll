package routes

import (
	controllers "AuthenticationService/controllers/Migrate"

	"github.com/gin-gonic/gin"
)

func InitMigrateRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/migrate/")
	route.GET("/dicom", controllers.DicomMigrateController())
	route.GET("/dicomone", controllers.DicomOneMigrateController())
}
