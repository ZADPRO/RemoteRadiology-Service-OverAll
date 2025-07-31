package routes

import (
	controllers "AuthenticationService/controllers/UserService"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"
)

func InitWellthgreenFormsRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/wellgreenforms")
	route.POST("/listPatientconsent", accesstoken.JWTMiddleware(), controllers.ListPatientConsentController())
	route.POST("/patientBrochure/list", accesstoken.JWTMiddleware(), controllers.ListPatientBrochureControllers())
	route.PATCH("/patientBrochure/update", accesstoken.JWTMiddleware(), controllers.UpdatePatientBrochureControllers())
	route.POST("/patientconsent/list", accesstoken.JWTMiddleware(), controllers.ListPatientConsentControllers())
	route.PATCH("/patientconsent/update", accesstoken.JWTMiddleware(), controllers.UpdatePatientConsentControllers())
	route.POST("/trainingmaterialguide/list", accesstoken.JWTMiddleware(), controllers.ListTrainingMaterialGuideControllers())
	route.PATCH("/trainingmaterialguide/update", accesstoken.JWTMiddleware(), controllers.UpdateTrainingMaterialGuideControllers())
	route.POST("/technicianconsent/list", accesstoken.JWTMiddleware(), controllers.ListTechnicianConsentControllers())
	route.PATCH("/technicianconsent/update", accesstoken.JWTMiddleware(), controllers.UpdateTechnicianConsentControllers())
}
