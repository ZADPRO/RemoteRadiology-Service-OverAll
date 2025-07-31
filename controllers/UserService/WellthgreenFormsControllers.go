package controllers

import (
	service "AuthenticationService/Service/UserService"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/UserService"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListPatientConsentController() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.ListPatientConsentReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ListPatientConsentService(dbConn, data)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"data": resVal,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func ListPatientBrochureControllers() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.ListPatientBrochureReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ListPatientBrochureService(dbConn, data, 1, 2)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":                 resVal.Status,
			"WGPatientBrochure":      resVal.WGPatientBrochure,
			"SCBrochureAccessStatus": resVal.SCBrochureAccessStatus,
			"SCPatientBrochure":      resVal.SCPatientBrochure,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func UpdatePatientBrochureControllers() gin.HandlerFunc {
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.UpdatePatientBroucherReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.UpdatePatientBrochureService(dbConn, data, int(idValue.(float64)), 1, 2, 36, 37)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(resVal, true, token),
			"token": token,
		})

	}
}

func ListPatientConsentControllers() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.ListPatientBrochureReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ListPatientBrochureService(dbConn, data, 3, 4)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":                 resVal.Status,
			"WGPatientBrochure":      resVal.WGPatientBrochure,
			"SCBrochureAccessStatus": resVal.SCBrochureAccessStatus,
			"SCPatientBrochure":      resVal.SCPatientBrochure,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func UpdatePatientConsentControllers() gin.HandlerFunc {
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.UpdatePatientBroucherReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.UpdatePatientBrochureService(dbConn, data, int(idValue.(float64)), 3, 4, 38, 39)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(resVal, true, token),
			"token": token,
		})

	}
}

func ListTrainingMaterialGuideControllers() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.ListPatientBrochureReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ListPatientBrochureService(dbConn, data, 5, 6)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":                 resVal.Status,
			"WGPatientBrochure":      resVal.WGPatientBrochure,
			"SCBrochureAccessStatus": resVal.SCBrochureAccessStatus,
			"SCPatientBrochure":      resVal.SCPatientBrochure,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func UpdateTrainingMaterialGuideControllers() gin.HandlerFunc {
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.UpdatePatientBroucherReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.UpdatePatientBrochureService(dbConn, data, int(idValue.(float64)), 5, 6, 40, 41)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(resVal, true, token),
			"token": token,
		})

	}
}

func ListTechnicianConsentControllers() gin.HandlerFunc {
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

		data, ok := helper.GetRequestBody[model.ListPatientBrochureReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.ListPatientBrochureService(dbConn, data, 7, 8)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		payload := map[string]interface{}{
			"status":                 resVal.Status,
			"WGPatientBrochure":      resVal.WGPatientBrochure,
			"SCBrochureAccessStatus": resVal.SCBrochureAccessStatus,
			"SCPatientBrochure":      resVal.SCPatientBrochure,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})

	}
}

func UpdateTechnicianConsentControllers() gin.HandlerFunc {
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.UpdatePatientBroucherReq](c, true)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.UpdatePatientBrochureService(dbConn, data, int(idValue.(float64)), 7, 8, 42, 43)
		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(resVal, true, token),
			"token": token,
		})

	}
}
