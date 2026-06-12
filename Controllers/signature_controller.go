package Controller

import (
	"net/http"
	"strconv"

	"encoding/json"

	"invoice-backend/Models/dto"
	hashapi "invoice-backend/Helper/HashAPI"

	"github.com/gin-gonic/gin"

	models "invoice-backend/Models"
    query "invoice-backend/Query"
	"fmt"
)

func GetSignatureAuthorities(c *gin.Context) {

	authorities, err := query.GetAllSignatureAuthorities()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//fmt.Println(authorities)

	c.JSON(http.StatusOK, authorities)
}

func CreateSignatureAuthority(c *gin.Context) {

	token := getToken(c)

	var packet dto.EncryptedPacket

	// 1. Bind encrypted packet
	if err := c.ShouldBindJSON(&packet); err != nil {

		c.JSON(http.StatusBadRequest,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Invalid Request Packet",
			}, true, token),
		)

		return
	}

	// 2. Decrypt request
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

	// 3. Convert decrypted data into struct
	var authority models.SignatureAuthority

	jsonBytes, _ := json.Marshal(decrypted)
	json.Unmarshal(jsonBytes, &authority)

	// 4. Save to DB
	err = query.CreateSignatureAuthority(&authority)

	if err != nil {

		//fmt.Println("CREATE SIGNATURE ERROR:", err)

		c.JSON(http.StatusInternalServerError,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": err.Error(),
			}, true, token),
		)

		return
	}

	// 5. Success response
	c.JSON(http.StatusOK,
		hashapi.Encrypt(gin.H{
			"status": true,
			"message": "Signature authority created successfully",
			"data": authority,
		}, true, token),
	)
}

func UpdateSignatureAuthority(c *gin.Context) {

	token := getToken(c)

	id, _ := strconv.Atoi(c.Param("id"))

	var packet dto.EncryptedPacket

	// 1. Bind encrypted packet
	if err := c.ShouldBindJSON(&packet); err != nil {

		c.JSON(http.StatusBadRequest,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Invalid Request Packet",
			}, true, token),
		)

		return
	}

	// 2. Decrypt request
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

	// 3. Convert decrypted data
	var authority models.SignatureAuthority

	jsonBytes, _ := json.Marshal(decrypted)
	json.Unmarshal(jsonBytes, &authority)

	authority.ID = id

	// 4. Update DB
	err = query.UpdateSignatureAuthority(authority)

	fmt.Printf("DEBUG: Received Authority Struct: %+v\n", authority)

	if err != nil {

		c.JSON(http.StatusInternalServerError,
			hashapi.Encrypt(gin.H{
				"status": false,
				"message": err.Error(),
			}, true, token),
		)

		return
	}

	// 5. Success response
	c.JSON(http.StatusOK,
		hashapi.Encrypt(gin.H{
			"status": true,
			"message": "Signature authority updated successfully",
			"data": authority,
		}, true, token),
	)
}

func DeleteSignatureAuthority(c *gin.Context) {

	token := getToken(c)

	id, _ := strconv.Atoi(c.Param("id"))

	err := query.DeleteSignatureAuthority(id)

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
			"message": "Signature authority deleted successfully",
		}, true, token),
	)
}