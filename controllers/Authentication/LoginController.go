package controllers

import (
	service "AuthenticationService/Service/Authentication"
	db "AuthenticationService/internal/DB"
	accesstoken "AuthenticationService/internal/Helper/AccessToken"
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	helper "AuthenticationService/internal/Helper/RequestHandler"
	model "AuthenticationService/internal/Model/Authentication"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginController() gin.HandlerFunc {

	return func(c *gin.Context) {

		var reqVal model.LoginReq

		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try Again " + err.Error(),
			})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.LoginServices(dbConn, reqVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		if resVal.Status {
			response["email"] = resVal.Email
			response["roleType"] = resVal.RoleType
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})

	}
}

func VerifyOTPController() gin.HandlerFunc {

	return func(c *gin.Context) {

		var reqVal model.VerifyReq

		if err := c.BindJSON(&reqVal); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": "Something went wrong, Try Again " + err.Error(),
			})
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		resVal := service.VerifyOTPService(dbConn, reqVal)

		response := gin.H{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		if resVal.Status {
			response["token"] = resVal.Token
			response["RoleType"] = resVal.RoleType
			response["PasswordStatus"] = resVal.PasswordStatus
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})

	}
}

func UserChangePasswordController() gin.HandlerFunc {
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

		reqVal, ok := helper.RequestHandler[model.UserChnagePasswordReq](c)
		if !ok {
			return
		}

		dbConn, sqlDB := db.InitDB()
		defer sqlDB.Close()

		// fmt.Println(reqVal)

		resVal := service.UserChangePasswordService(dbConn, *reqVal, idValue)

		payload := map[string]interface{}{
			"status":  resVal.Status,
			"message": resVal.Message,
		}

		token := accesstoken.CreateToken(idValue, roleIdValue)

		c.JSON(http.StatusOK, gin.H{
			"data":  hashapi.Encrypt(payload, true, token),
			"token": token,
		})
	}
}
