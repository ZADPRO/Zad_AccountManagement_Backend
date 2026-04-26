package Controller

import (
	"database/sql"
	"encoding/json"
	"invoice-backend/Helper/HashAPI"
	
	"invoice-backend/Models/dto" // Ensure you have this for EncryptedPacket
	"invoice-backend/Services"
	"net/http"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Helper to get token from context safely
func getToken(c *gin.Context) string {
	if val, exists := c.Get("raw_token"); exists {
		return val.(string)
	}
	return ""
}

func GetClientList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		clients, err := Services.GetAllClients(db)
		if err != nil {
			resp := hashapi.Encrypt(dto.BaseResponse{Status: false, Message: err.Error()}, true, token)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		// Wrap and Encrypt for Axios
		resp := hashapi.Encrypt(gin.H{"clients": clients}, true, token)
		c.JSON(http.StatusOK, resp)
	}
}

func CreateClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		adminIDRaw, _ := c.Get("user_id")
 		adminID := int(adminIDRaw.(float64))

		// 1. Bind the Encrypted Packet
		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Packet"}, true, token))
			return
		}

		// 2. Decrypt
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Decryption failed"}, true, token))
			return
		}

		// 3. Unmarshal into Request Struct
		// 3. Unmarshal into Request Struct
		var req dto.CreateClientRequest
		jsonBytes, _ := json.Marshal(decrypted)
		if err := json.Unmarshal(jsonBytes, &req); err != nil {
    		c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Payload structure mismatch"}, true, token))
    		return
		}

		// 4. Call Service (which handles Static DB Encryption)
		clientID, err := Services.CreateClient(db, req, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}
		fmt.Printf("DEBUG req: %+v\n", req)

		// 5. Encrypt Success Response
		success := hashapi.Encrypt(gin.H{
			"status":   true,
			"message":  "Client created successfully",
			"clientID": clientID,
		}, true, token)
		c.JSON(http.StatusOK, success)
	}
}

func UpdateClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		adminIDRaw, _ := c.Get("user_id")
 		adminID := int(adminIDRaw.(float64))

		idStr := c.Param("id")
		clientID, _ := strconv.Atoi(idStr)

		var packet dto.EncryptedPacket
		c.ShouldBindJSON(&packet)

		decrypted, _ := hashapi.Decrypt(packet.Data, token)
		var req dto.CreateClientRequest
		 jsonBytes, _ := json.Marshal(decrypted)
        if err := json.Unmarshal(jsonBytes, &req); err != nil {
            c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Payload format error"}, true, token))
            return
		}
		err := Services.UpdateClient(db, clientID, req, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "message": "Updated successfully"}, true, token))
	}
}

func GetClientByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		idParam := c.Param("id")
		clientID, _ := strconv.Atoi(idParam)

		client, err := Services.GetClientByID(db, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}

		// Encrypting the "data" key containing the client details
		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "data": client}, true, token))
	}
}

func DeleteClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		adminIDRaw, _ := c.Get("user_id")
		adminID := int(adminIDRaw.(float64))
		idStr := c.Param("id")
		clientID, _ := strconv.Atoi(idStr)

		err := Services.DeleteClient(db, clientID, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(dto.BaseResponse{Status: false, Message: err.Error()}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.BaseResponse{Status: true, Message: "Deactivated successfully"}, true, token))
	}
}