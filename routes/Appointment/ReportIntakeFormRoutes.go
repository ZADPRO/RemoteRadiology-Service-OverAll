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
	// route.POST("/answerreportIntake", accesstoken.JWTMiddleware(), controllers.AnswerReportIntakeController())
	// route.POST("/answerTechnicianIntake", accesstoken.JWTMiddleware(), controllers.AnswerTechnicianIntakeController())
	// route.POST("/answerPatientIntake", accesstoken.JWTMiddleware(), controllers.AnswerPatientIntakeController())
	// route.POST("/answerTextContent", accesstoken.JWTMiddleware(), controllers.AnswerTextContentController())
	route.POST("/addComments", accesstoken.JWTMiddleware(), controllers.AddCommentsController())
	// route.POST("/completeReport", accesstoken.JWTMiddleware(), controllers.CompleteReportController())
	route.POST("/submitReport", accesstoken.JWTMiddleware(), controllers.SubmitReportController())
	route.POST("/updateRemarks", accesstoken.JWTMiddleware(), controllers.UpdateRemarksController())
	route.POST("/uploadReportFormate", accesstoken.JWTMiddleware(), controllers.UploadReportFormateController())
	route.POST("/getReportFormate", accesstoken.JWTMiddleware(), controllers.GetReportFormateController())
	route.POST("/listremark", accesstoken.JWTMiddleware(), controllers.ListRemarkController())
	route.POST("/sendMail", accesstoken.JWTMiddleware(), controllers.SendMailReportController())
	route.POST("/downloadreport", accesstoken.JWTMiddleware(), controllers.DownloadReportService())
	// route.POST("/addAddendum", accesstoken.JWTMiddleware(), controllers.AddAddendumController())
	// route.POST("/sendMail", accesstoken.JWTMiddleware(), controllers.SendMailReportController())
}
