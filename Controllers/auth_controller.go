package Controller

import (
	"invoice-backend/DB"
    "database/sql"
	accesstoken "invoice-backend/Helper/AccessToken"
	hashapi "invoice-backend/Helper/HashAPI"
	"invoice-backend/Models/dto"
	"invoice-backend/Services"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)
// LoginHandler processes user authentication requests.
// It expects raw JSON, validates credentials, and returns an encrypted session token.
// Note: Errors are returned in plain JSON because the client hasn't established 
// an encrypted session context (token) yet.
func LoginHandler(c *gin.Context) {
    
    var loginReq dto.LoginRequest

    // Bind raw JSON request to struct. Use ShouldBindJSON for flexible error handling.
    if err := c.ShouldBindJSON(&loginReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request format"})
        return
    }

    // Logic: Database connection and credential verification via Services layer.
    db := DB.GetConn()
    userID, role, name, isFirst, err := Services.VerifyLogin(db, loginReq.Email, loginReq.Password)
    if err != nil {
        // Security: Keep error messages generic to prevent email enumeration.
        c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid email or password"})
        return
    }

   // Generate JWT: This token acts as the 'key' for all subsequent encrypted communication.
    token, err := accesstoken.CreateToken(userID) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Token error"})
        return
    }
	fmt.Println("TOKEN:", token)
    // Build the structured response.
    loginResp := dto.LoginResponse{
        BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
        Token:         token,
        Role:          role,
        Username:      name,
        IsFirstLogin:  isFirst,
        UserId:        userID,
    }

   // Response Strategy: Return the token in plain text so the frontend can store it,
	// but encrypt the sensitive user data using that same token.
    c.JSON(http.StatusOK, gin.H{
    "token": token,
    "data": hashapi.Encrypt(loginResp, true, token),
})
} 

// ChangePasswordHandler returns a middleware-ready HandlerFunc.
// This route uses "Dynamic Encryption," requiring the JWT from the header 
// to decrypt the incoming payload and encrypt the outgoing response.
func ChangePasswordHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Context Extraction
		// Extract token to use as the decryption/encryption key.
		token := getToken(c)

		// 2. Decryption Workflow
		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "Invalid encrypted packet",
			}, true, token))
			return
		}

		// Use the session token to decrypt the incoming request data.
		decryptedInterface, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "Decryption failed: identity could not be verified",
			}, true, token))
			return
		}

		// 3. Type Assertion & Validation
		// Ensure the decrypted data is a valid map and contains the required fields.
		data, ok := decryptedInterface.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "Data format error",
			}, true, token))
			return
		}

		newPassword, pOk := data["newPassword"].(string)
		if !pOk || newPassword == "" {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "New password is required",
			}, true, token))
			return
		}

		// 4. Authorization Check
		// We pull the UserID from the Gin Context (set by auth middleware) 
		// rather than the payload to prevent ID spoofing.
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "User session not found",
			}, true, token))
			return
		}
		// Handle potential type differences if UserID is stored as int vs float64 (common in JSON/JWT parsing).
		var finalID int
		switch v := userID.(type) {
		case int:
			finalID = v
		case float64:
			finalID = int(v)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid UserID type"})
			return
}
// 5. Execution
err = Services.ChangePassword(db, finalID, newPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
				"status":  false, 
				"message": "Failed to update password: " + err.Error(),
			}, true, token))
			return
		}

		// 6. Encrypted Success Response
		c.JSON(http.StatusOK, hashapi.Encrypt(dto.BaseResponse{
			Status:  true,
			Message: "Password updated successfully",
		}, true, token))
	}
}