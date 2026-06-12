package Controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	hashapi "invoice-backend/Helper/HashAPI"
	models "invoice-backend/Models"
	"invoice-backend/Models/dto"
	query "invoice-backend/Query"

	"github.com/gin-gonic/gin"
)

func GetCurrencies(c *gin.Context) {

	currencies, err := query.GetAllCurrencies()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, currencies)
}

func CreateCurrency(c *gin.Context) {

	token := getToken(c)

	var packet dto.EncryptedPacket

	if err := c.ShouldBindJSON(&packet); err != nil {

		c.JSON(http.StatusBadRequest,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Invalid Request Packet",
			}, true, token),
		)

		return
	}

	decrypted, err := hashapi.Decrypt(packet.Data, token)

	if err != nil {

		c.JSON(http.StatusUnauthorized,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Decryption failed",
			}, true, token),
		)

		return
	}

	var currency models.Currency

	jsonBytes, _ := json.Marshal(decrypted)
	json.Unmarshal(jsonBytes, &currency)

	err = query.CreateCurrency(&currency)

	if err != nil {

		fmt.Println("CREATE CURRENCY ERROR:", err)

		c.JSON(http.StatusInternalServerError,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": err.Error(),
			}, true, token),
		)

		return
	}

	c.JSON(http.StatusOK,
		hashapi.Encrypt(gin.H{
			"status": true,
			"message": "Currency created successfully",
			"data": currency,
		}, true, token),
	)
}

func UpdateCurrency(c *gin.Context) {

	token := getToken(c)

	id, _ := strconv.Atoi(c.Param("id"))

	var packet dto.EncryptedPacket

	if err := c.ShouldBindJSON(&packet); err != nil {

		c.JSON(http.StatusBadRequest,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Invalid Request Packet",
			}, true, token),
		)

		return
	}

	decrypted, err := hashapi.Decrypt(packet.Data, token)

	if err != nil {

		c.JSON(http.StatusUnauthorized,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Decryption failed",
			}, true, token),
		)

		return
	}

	var currency models.Currency

	jsonBytes, _ := json.Marshal(decrypted)
	json.Unmarshal(jsonBytes, &currency)

	currency.ID = id

	err = query.UpdateCurrency(currency)

	if err != nil {

		c.JSON(http.StatusInternalServerError,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": err.Error(),
			}, true, token),
		)

		return
	}

	c.JSON(http.StatusOK,
		hashapi.Encrypt(gin.H{
			"status": true,
			"message": "Currency updated successfully",
			"data": currency,
		}, true, token),
	)
}

func DeleteCurrency(c *gin.Context) {

	token := getToken(c)

	id, _ := strconv.Atoi(c.Param("id"))

	err := query.DeleteCurrency(id)

	if err != nil {

		c.JSON(http.StatusInternalServerError,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": err.Error(),
			}, true, token),
		)

		return
	}

	c.JSON(http.StatusOK,
		hashapi.Encrypt(gin.H{
			"status": true,
			"message": "Currency deleted successfully",
		}, true, token),
	)
}