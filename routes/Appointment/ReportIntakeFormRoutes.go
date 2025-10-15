package routes

import (
	controllers "AuthenticationService/controllers/Appointment"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"

)

func InitReportIntakeFormRoutes(router *gin.Engine) {
	route := router.Group("/api/v1/reportintakeform")
	route.POST("/checkaccess", accesstoken.JWTMiddleware(), controllers.CheckAccessController())
	route.POST("/assignreport", accesstoken.JWTMiddleware(), controllers.AssignGetReportController())
	route.POST("/addComments", accesstoken.JWTMiddleware(), controllers.AddCommentsController())
	route.POST("/submitReport", accesstoken.JWTMiddleware(), controllers.SubmitReportController())
	route.POST("/autosaveReport", accesstoken.JWTMiddleware(), controllers.AutosaveController())
	route.POST("/updateRemarks", accesstoken.JWTMiddleware(), controllers.UpdateRemarksController())
	route.POST("/uploadReportFormate", accesstoken.JWTMiddleware(), controllers.UploadReportFormateController())
	route.POST("/deleteReportFormate", accesstoken.JWTMiddleware(), controllers.DeleteReportFormateController())
	route.POST("/updateReportFormate", accesstoken.JWTMiddleware(), controllers.UpdateReportFormateController())
	route.POST("/getReportFormate", accesstoken.JWTMiddleware(), controllers.GetReportFormateController())
	route.POST("/listremark", accesstoken.JWTMiddleware(), controllers.ListRemarkController())
	route.POST("/sendMail", accesstoken.JWTMiddleware(), controllers.SendMailReportController())
	route.POST("/downloadreport", accesstoken.JWTMiddleware(), controllers.DownloadReportService())
	route.POST("/viewFiles", accesstoken.JWTMiddleware(), controllers.ViewReportService())
	route.POST("/addAddendum", accesstoken.JWTMiddleware(), controllers.AddAddendumController())
	route.POST("/listAllOldReport", accesstoken.JWTMiddleware(), controllers.ListOldReportController())
	route.POST("/addOldReport", accesstoken.JWTMiddleware(), controllers.PostOldReportUploadFileController())

	route.POST("/oldreportuploadurl", accesstoken.JWTMiddleware(), controllers.PostGenerateOldReportUploadURLController())

	route.POST("/deleteOldReport", accesstoken.JWTMiddleware(), controllers.DeleteOldReportController())
	route.POST("/insertSignature", accesstoken.JWTMiddleware(), controllers.InsertSignatureController())
}
