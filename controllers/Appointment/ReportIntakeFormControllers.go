package controllers

import (
	service "AuthenticationService/Service/Appointment"
	s3Service "AuthenticationService/Service/S3"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	logger "AuthenticationService/internal/Helper/Logger"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	timeZone "AuthenticationService/internal/Helper/TimeZone"
	helperView "AuthenticationService/internal/Helper/ViewFile"
	model "AuthenticationService/internal/Model/Appointment"
	s3config "AuthenticationService/internal/Storage/s3"
	query "AuthenticationService/query/Appointment"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func extractS3Key(s3URL string) string {
	prefix := "https://easeqt-health-archive.s3.us-east-2.amazonaws.com/"
	return strings.TrimPrefix(s3URL, prefix)
}

func CheckAccessController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.CheckAccessReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message, accessId, custId := service.CheckAccessService(dbConn, data, int(idValue.(float64)), int(roleIdValue.(float64)))

		payload := map[string]interface{}{
			"status":   status,
			"message":  message,
			"accessId": accessId,
			"custId":   custId,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AssignGetReportController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AssignGetReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message, IntakeFormData, TechnicianIntakeFormData, ReportIntakeFormData, ReportTextContentData, ReportHistoryData, ReportCommentsData, ReportAppointmentData, ReportFormateList, GetUserDetails, PatientUserDetails, EaseQTReportAccess, ScanCenterImg, ScancenterAddress, Addendum, oldReport, NASystemReportAccess, patientpublicprivate, PerformingProviderName, VerifyingProviderName, ListAllSignature, ReportPortalImpRecom, NAImpRecom := service.AssignGetReportService(dbConn, data, int(idValue.(float64)), int(roleIdValue.(float64)))

		payload := map[string]interface{}{
			"status":                   status,
			"message":                  message,
			"intakeFormData":           IntakeFormData,
			"technicianIntakeFormData": TechnicianIntakeFormData,
			"reportIntakeFormData":     ReportIntakeFormData,
			"reportTextContentData":    ReportTextContentData,
			"reportHistoryData":        ReportHistoryData,
			"reportCommentsData":       ReportCommentsData,
			"appointmentStatus":        ReportAppointmentData,
			"reportFormateList":        ReportFormateList,
			"userDeatils":              GetUserDetails,
			"patientDetails":           PatientUserDetails,
			"easeQTReportAccess":       EaseQTReportAccess,
			"ScanCenterImg":            ScanCenterImg,
			"ScancenterAddress":        ScancenterAddress,
			"Addendum":                 Addendum,
			"oldReport":                oldReport,
			"naSystemReportAccess":     NASystemReportAccess,
			"patientpublicprivate":     patientpublicprivate,
			"PerformingProviderName":   PerformingProviderName,
			"VerifyingProviderName":    VerifyingProviderName,
			"ListAllSignature":         ListAllSignature,
			"ReportPortalImpRecom":     ReportPortalImpRecom,
			"NAImpRecom":               NAImpRecom,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AnswerReportIntakeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AnswerReportIntakeReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AnswerReportIntakeService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AnswerTechnicianIntakeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AnswerReportIntakeReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AnswerTechnicianIntakeService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AnswerPatientIntakeController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AnswerReportIntakeReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AnswerPatientIntakeService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AnswerTextContentController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AnswerTextContentReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AnswerTextContentService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AddCommentsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AddCommentReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.AddCommentsService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func CompleteReportController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.CompleteReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.CompleteReportService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AutosaveController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AutoSubmitReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message, ReportIntake, TextContent, AppointmentDetails, EaseQTReportAccess, NASystemReportAccess, PerformingProviderName, VerifyingProviderName, ListAllSignature := service.AutosaveServicee(dbConn, data, int(idValue.(float64)), int(roleIdValue.(float64)))

		payload := map[string]interface{}{
			"status":                 status,
			"message":                message,
			"reportIntakeFormData":   ReportIntake,
			"reportTextContentData":  TextContent,
			"appointmentStatus":      AppointmentDetails,
			"easeQTReportAccess":     EaseQTReportAccess,
			"naSystemReportAccess":   NASystemReportAccess,
			"PerformingProviderName": PerformingProviderName,
			"VerifyingProviderName":  VerifyingProviderName,
			"ListAllSignature":       ListAllSignature,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func SubmitReportController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.SubmitReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.SubmitReportService(dbConn, data, int(idValue.(float64)), int(roleIdValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func UpdateRemarksController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.UpdateRemarkReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.UpdateRemarksService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func UploadReportFormateController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.UploadReportFormateReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		id, refUserCustId, status, message := service.UploadReportFormateService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"id":            id,
			"refUserCustId": refUserCustId,
			"status":        status,
			"message":       message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func DeleteReportFormateController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.DeleteReportFormateReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.DeleteReportFormateService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func UpdateReportFormateController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.UpdateReportFormateReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.UpdateReportFormateService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func GetReportFormateController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.GetReportFormateReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		TextData := service.GetReportFormateService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  true,
			"message": TextData,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func ListRemarkController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.ListRemarkReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		TextData := service.ListRemarkService(dbConn, data)

		payload := map[string]interface{}{
			"status":  true,
			"message": TextData,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func SendMailReportController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.SendMailReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.SendMailReportService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func DownloadReportService() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.DownloadReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		result := service.DownloadReportService(dbConn, data)

		payload := map[string]interface{}{
			"status": true,
			"data":   result,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

// func ViewReportService() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		logger := logger.InitLogger()

// 		idValue, idExists := c.Get("id")
// 		roleIdValue, roleIdExists := c.Get("roleId")

// 		if !idExists || !roleIdExists {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"status":  false,
// 				"message": "User ID, RoleID not found in request context.",
// 			})
// 			return
// 		}

// 		data, ok := helper.GetRequestBody[model.ViewReportReq](c, true)
// 		if !ok {
// 			return
// 		}

// 		// Default payload
// 		payload := map[string]interface{}{
// 			"status": false,
// 			"data": map[string]interface{}{
// 				"base64Data":  "",
// 				"contentType": "",
// 			},
// 		}

// 		// Try to load the file
// 		ViewFiles, viewErr := helperView.ViewFile("./Assets/Files/" + data.FileName)
// 		if viewErr != nil {
// 			logger.Printf("Failed to read DrivingLicense file: %v", viewErr)
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"status":  false,
// 				"message": "Failed to load file",
// 			})
// 			return
// 		}

// 		// Update payload on success
// 		payload["status"] = true
// 		payload["data"] = map[string]interface{}{
// 			"base64Data":  ViewFiles.Base64Data,
// 			"contentType": ViewFiles.ContentType,
// 		}

// 		// Create token
// 		token := accesstoken.CreateToken(idValue, roleIdValue)

// 		// Respond
// 		c.JSON(http.StatusOK, gin.H{
// 			"data":  hashapi.Encrypt(payload, true, token),
// 			"token": token,
// 		})
// 	}
// }

func ViewReportService() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.InitLogger()

		// Get user ID and role from context
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")
		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID, RoleID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.ViewReportReq](c, true)
		if !ok {
			return
		}

		payload := map[string]interface{}{
			"status": false,
			"data": map[string]interface{}{
				"base64Data":  "",
				"contentType": "",
			},
		}

		if strings.HasPrefix(data.FileName, "https://easeqt-health-archive.s3") {
			// Always generate presigned URL for S3 files
			key := extractS3Key(data.FileName)
			presignedURL, err := s3Service.GeneratePresignGetURL(context.Background(), key, 10*time.Minute)
			if err != nil {
				logger.Printf("Failed to generate presigned URL: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to generate file URL",
				})
				return
			}

			// Print the presigned URL to console
			fmt.Println("Generated Presigned URL:", presignedURL)
			// or using logger
			logger.Printf("Generated Presigned URL: %s", presignedURL)

			payload["status"] = true
			payload["data"] = map[string]interface{}{
				"base64Data":  presignedURL,
				"contentType": "url",
			}
		} else {
			// fmt.Println("\n\nElse Block =>>>>>>>>> \n\n")
			// Local file — read and encode
			ViewFiles, viewErr := helperView.ViewFile("./Assets/Files/" + data.FileName)
			if viewErr != nil {
				logger.Printf("Failed to read local file: %v", viewErr)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to load file",
				})
				return
			}

			payload["status"] = true
			payload["data"] = map[string]interface{}{
				"base64Data":  ViewFiles.Base64Data,
				"contentType": ViewFiles.ContentType,
			}
		}

		// Create token
		token := accesstoken.CreateToken(idValue, roleIdValue)

		fmt.Printf("\n\npayload: %+v\n\n", payload)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func AddAddendumController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AddAddendumReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message, resVal := service.AddAddendumService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    resVal,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func ListOldReportController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.ListOldReportReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message, resVal := service.ListOldReportService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    resVal,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostOldReportUploadFileController() gin.HandlerFunc {
	return func(c *gin.Context) {

		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return // Stop processing
		}

		uploadPath := "./Assets/Files/"

		log := logger.InitLogger()

		file, err := c.FormFile("file")
		patientId := c.PostForm("patientId")
		categoryId := c.PostForm("categoryId")
		appointmentId := c.PostForm("appointmentId")
		if err != nil || patientId == "" || categoryId == "" || appointmentId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Error retrieving profile image from request: " + err.Error(),
			})
			return
		}

		maxFileSize := int64(10 * 1024 * 1024) // 10 MB
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": fmt.Sprintf("Profile image size exceeds the limit of %d MB", maxFileSize/(1024*1024)),
			})
			return
		}

		ext := filepath.Ext(file.Filename)
		allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf"}
		isAllowed := false

		for _, allowedExt := range allowedExts {
			if strings.ToLower(ext) == allowedExt {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid profile image file type. Only JPG, JPEG, PNG are allowed.",
			})
			return
		}

		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),                           // Generate a random UUID
			timeZone.GetTimeWithFormate("20060102150405"), // Add timestamp (YYYYMMDDHHMMSS)
			ext) // Keep original file extension
		destinationPath := filepath.Join(uploadPath, uniqueFilename)

		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			log.Printf("Error creating upload directory '%s': %v\n", uploadPath, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Server error: Could not prepare image storage.",
			})
			return
		}

		if err := c.SaveUploadedFile(file, destinationPath); err != nil {
			log.Printf("Error saving uploaded file to '%s': %v\n", destinationPath, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Server error: Could not save profile image.",
			})
			return
		}

		log.Printf("Successfully uploaded image: %s\n", destinationPath)

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		InsertFileErr := dbConn.Exec(
			query.AddOldReportSQL,
			patientId,
			appointmentId,
			categoryId,
			uniqueFilename,
			timeZone.GetPacificTime(),
			idValue,
		).Error

		if InsertFileErr != nil {
			log.Printf("ERROR: Failed to Insert Old Report: %v\n", InsertFileErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Something went wrong, Try Again",
			})
			return
		}

		payload := map[string]interface{}{
			"status":      true,
			"message":     "File uploaded successfully!",
			"fileName":    uniqueFilename,
			"oldFilename": file.Filename,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostGenerateOldReportUploadURLController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "User ID or Role ID not found.",
			})
			return
		}

		var req struct {
			FileName      string `json:"fileName"`
			PatientId     string `json:"patientId"`
			CategoryId    string `json:"categoryId"`
			AppointmentId string `json:"appointmentId"`
		}

		if err := c.BindJSON(&req); err != nil || req.FileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid request. File name missing.",
			})
			return
		}

		ext := filepath.Ext(req.FileName)
		allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf"}
		isAllowed := false
		for _, a := range allowedExts {
			if strings.EqualFold(ext, a) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Invalid file type. Only JPG, JPEG, PNG, PDF allowed.",
			})
			return
		}

		uniqueFilename := fmt.Sprintf("%s_%s%s",
			uuid.New().String(),
			timeZone.GetTimeWithFormate("20060102150405"),
			ext,
		)

		s3Key := fmt.Sprintf("oldReportsPatient/%s", uniqueFilename)
		ctx := c.Request.Context()

		s3Client, err := s3config.New(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to initialize S3 client.",
			})
			return
		}

		uploadURL, err := s3Client.PresignPut(ctx, s3Key, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to generate upload URL.",
			})
			return
		}

		viewURL, err := s3Client.PresignGet(ctx, s3Key, 10*time.Minute)
		if err != nil {
			viewURL = ""
		}

		cleanViewURL := viewURL
		if idx := strings.Index(viewURL, "?"); idx != -1 {
			cleanViewURL = viewURL[:idx]
		}
		// ✅ DB connection
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		var count int64
		checkErr := dbConn.Raw(`
				SELECT COUNT(1) 
				FROM notes."refOldReport"
				WHERE "refUserId" = ? AND "refAppointmentId" = ? AND "refORCategoryId" = ? AND "refORFilename" = ?`,
			req.PatientId, req.AppointmentId, req.CategoryId, cleanViewURL,
		).Scan(&count).Error

		if checkErr != nil {
			log.Printf("ERROR checking existing old report: %v", checkErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to check existing report.",
			})
			return
		}

		// If not exists, insert into DB
		if count == 0 {
			insertErr := dbConn.Exec(
				query.AddOldReportSQL,
				req.PatientId,
				req.AppointmentId,
				req.CategoryId,
				cleanViewURL, // use clean URL instead of filename
				timeZone.GetPacificTime(),
				idValue,
			).Error

			if insertErr != nil {
				log.Printf("ERROR inserting old report: %v", insertErr)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "Failed to save report in DB.",
				})
				return
			}
		}

		payload := map[string]interface{}{
			"status":        true,
			"message":       "Old Report Upload URL generated successfully!",
			"uploadURL":     uploadURL,
			"viewURL":       viewURL,
			"s3Key":         s3Key,
			"fileName":      uniqueFilename,
			"oldFileName":   req.FileName,
			"patientId":     req.PatientId,
			"categoryId":    req.CategoryId,
			"appointmentId": req.AppointmentId,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func DeleteOldReportController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.DeleteOldReportModel](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.DeleteOldReportService(dbConn, data)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func InsertSignatureController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idValue, idExists := c.Get("id")
		roleIdValue, roleIdExists := c.Get("roleId")

		if !idExists || !roleIdExists {
			// Handle error: ID is missing from context (e.g., middleware didn't set it)
			c.JSON(http.StatusUnauthorized, gin.H{ // Or StatusInternalServerError depending on why it's missing
				"status":  false,
				"message": "User ID, RoleID, Branch ID not found in request context.",
			})
			return
		}

		data, ok := helper.GetRequestBody[model.AddSignatureReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message, list := service.InsertSignatureService(dbConn, data)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    list,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
