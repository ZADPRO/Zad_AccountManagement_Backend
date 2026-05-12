package Controller

import (
    "database/sql"
    "encoding/json" 
    "net/http"
    
    "invoice-backend/Models/dto"
    "invoice-backend/Services"
	"invoice-backend/Models/responses"
    "invoice-backend/Helper/HashAPI"  
	"strconv"
    
    "github.com/gin-gonic/gin"
)

// --- BANKING MANAGEMENT ---

// CreateBanking handles the creation of new bank account details.
// It decrypts the incoming banking packet and associates it with the current admin's ID.
func CreateBanking(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

        // Authentication: Extract Admin ID from the session context
		adminIDRaw, _ := c.Get("user_id")
		adminID := int(adminIDRaw.(float64))

        // 1. Packet Binding
		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status": false, "message": "Invalid Packet",
			}, true, token))
			return
		}

		// 2. Security Handshake (Decryption)
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status": false, "message": "Decryption failed",
			}, true, token))
			return
		}

		// 3. Transformation: Map the generic decrypted interface to a typed Banking DTO
		var req dto.SaveBankingRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		// 4. Persistence: Pass to service layer for SQL insertion
		id, err := Services.CreateBankingDetails(db, req, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
				"status": false, "message": err.Error(),
			}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{
			"status":    true,
			"message":   "Bank created successfully",
			"detailsId": id,
		}, true, token))
	}
} 

// UpdateBanking modifies existing bank records.
// Logic Note: We pull the ID from the URL path to ensure the user isn't 
// attempting to overwrite a different record by spoofing the JSON ID.
func UpdateBanking(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		adminIDRaw, _ := c.Get("user_id")
		adminID := int(adminIDRaw.(float64))

		// Get ID from RESTful path: /banking/:id
		idStr := c.Param("id")
		detailsID, _ := strconv.Atoi(idStr)

		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status": false, "message": "Invalid Packet",
			}, true, token))
			return
		}

		// Decrypt
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status": false, "message": "Decryption failed",
			}, true, token))
			return
		}

		var req dto.SaveBankingRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		// DATA INTEGRITY: Override any ID in the body with the ID from the URL path
		req.DetailsID = detailsID

		// 🔁 Call UPDATE service
		id, err := Services.UpdateBankingDetails(db, req, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
				"status": false, "message": err.Error(),
			}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{
			"status":    true,
			"message":   "Bank updated successfully",
			"detailsId": id,
		}, true, token))
	}
}


// GetBankingInfo retrieves all accounts associated with the current user.
func GetBankingInfo(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c)
        
        // Extract user_id set by your Auth Middleware
        val, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
            return
        }
        userID := int(val.(float64))

        data, err := Services.GetAllBankingDetails(db, userID)
        if err != nil {
            // Return success: true but empty array so React map() doesn't break
            c.JSON(http.StatusOK, hashapi.Encrypt(responses.BankingDetailsResponse{
                BaseResponse: dto.BaseResponse{Status: true, Message: "No accounts found"},
                Data:         []responses.BankingDetailsData{}, 
            }, true, token))
            return
        }

        c.JSON(http.StatusOK, hashapi.Encrypt(responses.BankingDetailsResponse{
            BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
            Data:         data,
        }, true, token))
    }
}

// DeleteBankingInfo performs a soft delete on banking records.
func DeleteBankingInfo(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c)
        
        // 1. Get DetailsID from URL
        detailsIDStr := c.Param("id")
        detailsID, _ := strconv.Atoi(detailsIDStr)

        // 2. Get UserID from Context
        val, _ := c.Get("user_id")
        userID := int(val.(float64))

        // 3. Call Service
        err := Services.SoftDeleteBankingDetails(db, detailsID, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
                "status": false, 
                "message": "Failed to delete account",
            }, true, token))
            return
        }

        c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{
            "status": true,
            "message": "Bank account deleted successfully",
        }, true, token))
    }
}


// --- CUSTOM FIELDS MANAGEMENT ---

// CreateCustomField adds dynamic inputs to invoices/clients.
// It returns the newly created field data so the frontend can update 
// its state immediately without a full page refresh.
func CreateCustomField(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c)
        var packet dto.EncryptedPacket

        if err := c.ShouldBindJSON(&packet); err != nil {
            c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Request Packet"}, true, token))
            return
        }

        decrypted, err := hashapi.Decrypt(packet.Data, token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Security Handshake Failed"}, true, token))
            return
        }

        var req dto.CreateCustomFieldRequest
        jsonBytes, _ := json.Marshal(decrypted)
        json.Unmarshal(jsonBytes, &req)

        val, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "User context missing"}, true, token))
            return
        }
        userID := int(val.(float64))

        id, err := Services.AddCustomField(db, req, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Failed to save field"}, true, token))
            return
        }

        c.JSON(http.StatusCreated, hashapi.Encrypt(gin.H{
            "status":     true,
            "message":    "Custom field created successfully",
            "fieldId":    id,
            "fieldLabel": req.FieldLabel, 
            "fieldType":  req.FieldType,  
            "isRequired": req.IsRequired, 
        }, true, token))
    }
}

// GetCustomFieldList fetches all configured custom fields.
func GetCustomFieldList(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c)

        // 1. Get plain text data from Service
        data, err := Services.FetchAllCustomFields(db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Could not retrieve fields"}, true, token))
            return
        }

        // 2. Prepare the response object
        resp := responses.CustomFieldListResponse{
            BaseResponse: dto.BaseResponse{
                Status:  true,
                Message: "Fields fetched successfully",
            },
            Fields: data,
        }

        // 3. ENCRYPT the whole list for the network handshake
        c.JSON(http.StatusOK, hashapi.Encrypt(resp, true, token))
    }
}

// DeleteCustomField handles the ID and encrypts the status response
func DeleteCustomField(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c)
        
        idStr := c.Param("id")
        fieldID, err := strconv.Atoi(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid ID"}, true, token))
            return
        } 
		  // ✅ Get user ID from middleware
        val, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
                "status": false,
                "message": "User not found",
            }, true, token))
            return
        }

        userID := int(val.(float64))


        // Call plain-text service
        err = Services.RemoveCustomField(db, fieldID,userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Delete failed"}, true, token))
            return
        }

        // Encrypt success message
        c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "message": "Deleted"}, true, token))
    }
} 

