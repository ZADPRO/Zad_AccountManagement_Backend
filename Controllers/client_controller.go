package Controller

import (
	"database/sql"
	"encoding/json"
	"invoice-backend/Helper/HashAPI"
	
	"invoice-backend/Models/dto" 
	"invoice-backend/Services"
	"net/http"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
)

// getToken retrieves the raw JWT string from the Gin context.
// This token is used as the symmetric key for encrypting and decrypting 
// request/response payloads for the current session.
func getToken(c *gin.Context) string {
	if val, exists := c.Get("raw_token"); exists {
		return val.(string)
	}
	return ""
}

// GetClientList fetches all clients from the database.
// The resulting list is wrapped in a JSON object and encrypted before transmission.
func GetClientList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		clients, err := Services.GetAllClients(db)
		if err != nil {
			// Encrypt error messages to maintain consistent communication protocol
			resp := hashapi.Encrypt(dto.BaseResponse{Status: false, Message: err.Error()}, true, token)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		// Success response: { "clients": [...] } -> Encrypted
		resp := hashapi.Encrypt(gin.H{"clients": clients}, true, token)
		c.JSON(http.StatusOK, resp)
	}
}

// CreateClient handles the registration of a new client.
// It decrypts the incoming packet, transforms the data into a CreateClientRequest struct,
// and passes it to the service layer.
func CreateClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		// Extract Admin ID from middleware context (authentication check)
		adminIDRaw, _ := c.Get("user_id")
 		adminID := int(adminIDRaw.(float64))

		// 1. Bind the Encrypted Packet (the outer envelope)
		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Packet"}, true, token))
			return
		}

		// 2. Decrypt the payload using the session token
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Decryption failed"}, true, token))
			return
		}

		// 3. Unmarshal the generic decrypted interface into the specific DTO
		// We marshal back to JSON then Unmarshal to leverage struct tags and validation
		var req dto.CreateClientRequest
		jsonBytes, _ := json.Marshal(decrypted)
		if err := json.Unmarshal(jsonBytes, &req); err != nil {
    		c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Payload structure mismatch"}, true, token))
    		return
		}

		// 4. Persistence: The service layer handles DB-level encryption if required
		clientID, err := Services.CreateClient(db, req, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}
		fmt.Printf("DEBUG req: %+v\n", req)

		// 5. Encrypt and return the success notification
		success := hashapi.Encrypt(gin.H{
			"status":   true,
			"message":  "Client created successfully",
			"clientID": clientID,
		}, true, token)
		c.JSON(http.StatusOK, success)
	}
}

// UpdateClient modifies an existing client record identified by the 'id' URL parameter.
func UpdateClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		adminIDRaw, _ := c.Get("user_id")
 		adminID := int(adminIDRaw.(float64))

		// Parse the client ID from the route path: /clients/:id
		idStr := c.Param("id")
		clientID, _ := strconv.Atoi(idStr)

		// Decryption and Binding logic
		var packet dto.EncryptedPacket
		c.ShouldBindJSON(&packet)// Note: Error check recommended here for production

		decrypted, _ := hashapi.Decrypt(packet.Data, token)
		var req dto.CreateClientRequest
		 jsonBytes, _ := json.Marshal(decrypted)
        if err := json.Unmarshal(jsonBytes, &req); err != nil {
            c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Payload format error"}, true, token))
            return
		}

		// Call service to update record
		err := Services.UpdateClient(db, clientID, req, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "message": "Updated successfully"}, true, token))
	}
}

// GetClientByID retrieves a single client's details.
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

		// Wrap client data in a status object and encrypt for secure transit
		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "data": client}, true, token))
	}
}

// DeleteClient performs a 'soft delete' or deactivation of a client.
// It requires the Admin ID from the session to log who performed the action.
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