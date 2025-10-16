package controllers

import (
	service "AuthenticationService/Service/UserService"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/UserService"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostCheckPatientController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.PatientCheckReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PostCheckPatientService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func PostAddPatientController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.RegisterNewPatientReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PostCheckPatientService(dbConn, model.PatientCheckReq{
			PatientId: data.PatientId,
			EmailId:   data.EmailId,
			PhoneNo:   data.PhoneNo,
		}, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		if !status {
			payload := map[string]interface{}{
				"status":  status,
				"message": message,
			}

			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
		} else {
			status, message := service.PostPatientService(dbConn, *data, int(idValue.(float64)))
			token := accesstoken.CreateToken(idValue, roleIdValue)

			payload := map[string]interface{}{
				"status":  status,
				"message": message,
			}

			c.JSON(http.StatusOK, gin.H{
				"data":  hashapi.Encrypt(payload, true, token),
				"token": token,
			})
		}

	}
}

func PatchUpdatePatientController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.UpdatePatientReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		fmt.Println("---------->HEllo", data)

		status, message := service.PatchPatientService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostCreatePatientController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.CreateAppointmentPatientReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PostCreatePatientService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostSendMailPatientController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.CreateAppointmentPatientReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PostSendMailPatientService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func PostCancelResheduleAppointmentController() gin.HandlerFunc {
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

		data, ok := helper.RequestHandler[model.CancelResheduleAppointmentReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		status, message := service.PostCancelResheduleAppointmentService(dbConn, *data, int(idValue.(float64)))
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
