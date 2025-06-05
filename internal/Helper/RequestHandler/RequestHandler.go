package helper

import (
	hashapi "AuthenticationService/internal/Helper/HashAPI"
	model "AuthenticationService/internal/Model/UserService"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func RequestHandler[T any](c *gin.Context) (*T, bool) {
	// Extract token from context
	tokenVal, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Token not found in context.",
		})
		return nil, false
	}

	// Bind encrypted body
	var encryptedData model.ReqVal
	if err := c.BindJSON(&encryptedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
		})
		return nil, false
	}

	if len(encryptedData.EncryptedData) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid encrypted data format",
		})
		return nil, false
	}

	// Decrypt
	decryptedInterface, err := hashapi.Decrypt(encryptedData.EncryptedData, tokenVal.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Decryption failed: " + err.Error(),
		})
		return nil, false
	}

	// fmt.Println("---> Decrypted MapData-----", decryptedInterface)

	// Validate decrypted structure
	mapData, ok := decryptedInterface.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Invalid decrypted format",
		})
		return nil, false
	}

	// fmt.Println("---> Decrypted MapData", mapData)
	fmt.Println("--- Decoded Struct Data:", mapData)

	var data T
	if err := mapstructure.Decode(mapData, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to decode decrypted data: " + err.Error(),
		})
		return nil, false
	}

	return &data, true
}
