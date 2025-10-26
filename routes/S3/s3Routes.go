package s3Routes

import (
	s3Controller "AuthenticationService/controllers/S3"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"

	"github.com/gin-gonic/gin"

)

func InitS3Routes(router *gin.Engine) {
	route := router.Group("/api/v1/storage/s3")
	{
		route.GET("/aws-s3-create", s3Controller.S3GeneratePresignPutController())
		route.GET("/aws-s3-read", s3Controller.S3GeneratePresignGetController())
		route.GET("/aws-s3-presign", s3Controller.S3GetFileController())
		// route.GET("/check", s3Controller.AckCheckController())
		route.POST("/final-report-upload", accesstoken.JWTMiddleware(), s3Controller.S3FinalReportUploadController())

	}
}
