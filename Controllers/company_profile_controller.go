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

func GetCompanyProfiles(c *gin.Context) {

	profiles, err := query.GetAllCompanyProfiles()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

func CreateCompanyProfile(c *gin.Context) {

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

	var profile models.CompanyProfile

	jsonBytes, _ := json.Marshal(decrypted)
	json.Unmarshal(jsonBytes, &profile)

	err = query.CreateCompanyProfile(&profile)

	if err != nil {

		fmt.Println("CREATE COMPANY PROFILE ERROR:", err)

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
			"message": "Company profile created successfully",
			"data": profile,
		}, true, token),
	)
}

func UpdateCompanyProfile(c *gin.Context) {

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

	var profile models.CompanyProfile

	jsonBytes, _ := json.Marshal(decrypted)
	json.Unmarshal(jsonBytes, &profile)

	profile.ID = id

	err = query.UpdateCompanyProfile(profile)

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
			"message": "Company profile updated successfully",
			"data": profile,
		}, true, token),
	)
}

func DeleteCompanyProfile(c *gin.Context) {

	token := getToken(c)

	id, _ := strconv.Atoi(c.Param("id"))

	err := query.DeleteCompanyProfile(id)

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
			"message": "Company profile deleted successfully",
		}, true, token),
	)
}