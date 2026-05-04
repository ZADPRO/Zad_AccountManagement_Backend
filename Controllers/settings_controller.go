package Controller

import (
    "database/sql"
    "encoding/json" // Added missing import
    "net/http"
    
    "invoice-backend/Models/dto"
    "invoice-backend/Services"
	"invoice-backend/Models/responses"
    "invoice-backend/Helper/HashAPI" // Ensure this matches your folder structure 
	"strconv"
    
    "github.com/gin-gonic/gin"
)

// CreateOrUpdateBanking handles the POST request from the "Add Bank Details" Modal
func CreateBanking(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		adminIDRaw, _ := c.Get("user_id")
		adminID := int(adminIDRaw.(float64))

		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status": false, "message": "Invalid Packet",
			}, true, token))
			return
		}

		// 🔓 Decrypt
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status": false, "message": "Decryption failed",
			}, true, token))
			return
		}

		// 🔁 Map to DTO
		var req dto.SaveBankingRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		// 🆕 Call INSERT service
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

func UpdateBanking(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		adminIDRaw, _ := c.Get("user_id")
		adminID := int(adminIDRaw.(float64))

		// 📌 Get ID from URL
		idStr := c.Param("id")
		detailsID, _ := strconv.Atoi(idStr)

		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status": false, "message": "Invalid Packet",
			}, true, token))
			return
		}

		// 🔓 Decrypt
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

		// 🔥 IMPORTANT: force ID from URL (don’t trust frontend blindly)
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



// GetBankingInfo retrieves the data for the Settings Page
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

// CreateCustomField handles the network handshake and calls the plain-text service
func CreateCustomField(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c) // Get session token for the handshake
        var packet dto.EncryptedPacket

        // 1. Bind the encrypted packet coming from React
        if err := c.ShouldBindJSON(&packet); err != nil {
            c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Request Packet"}, true, token))
            return
        }

        // 2. Decrypt the network packet into plain JSON
        decrypted, err := hashapi.Decrypt(packet.Data, token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Security Handshake Failed"}, true, token))
            return
        }

        // 3. Map Decrypted Data to the DTO
        var req dto.CreateCustomFieldRequest
        jsonBytes, _ := json.Marshal(decrypted)
        json.Unmarshal(jsonBytes, &req)

        // 4. Get User ID from Middleware (for audit log)
        // 4. Get User ID safely
val, exists := c.Get("user_id")
if !exists {
    // If user_id isn't in context, the middleware failed or is missing
    c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "User context missing"}, true, token))
    return
}

// Use float64 assertion if your JWT library stores numbers that way
userID := int(val.(float64))
        // 5. Call Service (Service stores it as PLAIN TEXT)
        id, err := Services.AddCustomField(db, req, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Failed to save field"}, true, token))
            return
        }

        // 6. ENCRYPT the response so React can decrypt it
        c.JSON(http.StatusCreated, hashapi.Encrypt(gin.H{
            "status":  true,
            "message": "Custom field created successfully",
            "fieldId": id,
        }, true, token))
    }
}

// GetCustomFieldList retrieves plain data from Service and encrypts it for the Wire
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

