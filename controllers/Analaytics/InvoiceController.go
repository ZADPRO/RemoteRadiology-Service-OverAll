package controllers

import (
	service "AuthenticationService/Service/Analaytics"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Analaytics"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAmountController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		Status, AmountModel, ScancenterData, UserData := service.GetAmountService(dbConn)

		//Load the Payload
		payload := map[string]interface{}{
			"status":         Status,
			"AmountModel":    AmountModel,
			"scancenterData": ScancenterData,
			"userData":       UserData,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func UpdateAmountController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.AmountModel](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		Status, message := service.UpdateAmountService(dbConn, data)

		//Load the Payload
		payload := map[string]interface{}{
			"status":  Status,
			"message": message,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func GetInvoiceDataController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.GetInvoiceDataReq](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		response := service.GetInvoiceDataService(dbConn, data)

		var payload map[string]interface{}

		//Load the Payload
		switch data.Type {
		case 1:
			payload = map[string]interface{}{
				"status":                true,
				"refTASform":            response.AmountModel[0].TASform,
				"refTADaform":           response.AmountModel[0].TADaform,
				"refTADbform":           response.AmountModel[0].TADbform,
				"refTADcform":           response.AmountModel[0].TADcform,
				"refTAXform":            response.AmountModel[0].TAXform,
				"refTAEditform":         response.AmountModel[0].TAEditform,
				"refTADScribeTotalcase": response.AmountModel[0].TADScribeTotalcase,
				"ScancenterData":        response.ScanCenterModel,
				"ScanCenterCount":       response.GetCountScanCenterMonthModel,
				"UserData":              nil,
				"UserCount":             nil,
			}

		case 2:
			payload = map[string]interface{}{
				"status":                true,
				"refTASform":            response.AmountModel[0].TASform,
				"refTADaform":           response.AmountModel[0].TADaform,
				"refTADbform":           response.AmountModel[0].TADbform,
				"refTADcform":           response.AmountModel[0].TADcform,
				"refTAXform":            response.AmountModel[0].TAXform,
				"refTAEditform":         response.AmountModel[0].TAEditform,
				"refTADScribeTotalcase": response.AmountModel[0].TADScribeTotalcase,
				"ScancenterData":        nil,
				"ScanCenterCount":       nil,
				"UserData":              response.GetUserModel,
				"UserCount":             response.AdminOverallScanIndicatesAnalayticsModel,
			}

		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func GenerateInvoiceDataController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.GenerateInvoiceReq](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		status, message := service.GenerateInvoiceDataService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":  status,
			"message": message,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func GetInvoiceHistoryController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.GetInvoiceHistoryReq](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		invoiceHistory, invoiceHistoryTakenDate := service.GetInvoiceHistoryService(dbConn, data, int(idValue.(float64)))

		payload := map[string]interface{}{
			"status":                  true,
			"invoiceHistory":          invoiceHistory,
			"invoiceHistoryTakenDate": invoiceHistoryTakenDate,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}

func GetInvoiceOverAllHistoryController() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Gathering the Datas From the Token
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

		//Request Should Be Encrypt
		data, ok := helper.GetRequestBody[model.GetInvoiceOverAllHistoryReq](c, true)
		if !ok {
			return
		}

		//DB Connections Intitied
		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		//Request Pass to the Service and Get the Retrun Value
		resVal := service.GetInvoiceOverAllHistoryService(dbConn, data, int(roleIdValue.(float64)))

		payload := map[string]interface{}{
			"status":         true,
			"invoiceHistory": resVal,
		}

		//Create a tokens
		token := accesstoken.CreateToken(idValue, roleIdValue)

		//Send a Reponse
		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
