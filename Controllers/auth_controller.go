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

func LoginHandler(c *gin.Context) {
    // 1. Change to a regular LoginRequest struct
    var loginReq dto.LoginRequest

    // 2. Bind regular JSON (This fixes the 401/400 errors)
    if err := c.ShouldBindJSON(&loginReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request format"})
        return
    }

    // 3. Verify Credentials (Directly using the bound fields)
    db := DB.GetConn()
    userID, role, name, isFirst, err := Services.VerifyLogin(db, loginReq.Email, loginReq.Password)
    if err != nil {
        // Return plain JSON error because frontend doesn't have a token to decrypt this yet
        c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid email or password"})
        return
    }

    // 4. Create Token
    token, err := accesstoken.CreateToken(userID) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Token error"})
        return
    }
	fmt.Println("TOKEN:", token)
    // 5. Prepare Response
    loginResp := dto.LoginResponse{
        BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
        Token:         token,
        Role:          role,
        Username:      name,
        IsFirstLogin:  isFirst,
        UserId:        userID,
    }

    // 🔐 6. Encrypt only the SUCCESS response
    // The frontend interceptor will use the token to decrypt this
    c.JSON(http.StatusOK, gin.H{
    "token": token,
    "data": hashapi.Encrypt(loginResp, true, token),
})
}
func ChangePasswordHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the raw JWT token for dynamic decryption/encryption
		token := getToken(c)

		// 2. Bind the encrypted packet from Axios
		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "Invalid encrypted packet",
			}, true, token))
			return
		}

		// 3. DYNAMIC DECRYPTION: Unlock the payload using the session token
		decryptedInterface, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "Decryption failed: identity could not be verified",
			}, true, token))
			return
		}

		// 4. Extract data from decrypted map
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

		// 5. SECURE IDENTITY: Pull the UserID directly from the JWT context
		// This prevents User A from changing User B's password
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status":  false,
				"message": "User session not found",
			}, true, token))
			return
		}
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

err = Services.ChangePassword(db, finalID, newPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
				"status":  false, 
				"message": "Failed to update password: " + err.Error(),
			}, true, token))
			return
		}

		// 7. DYNAMIC ENCRYPTION: Return success message locked with the token
		c.JSON(http.StatusOK, hashapi.Encrypt(dto.BaseResponse{
			Status:  true,
			Message: "Password updated successfully",
		}, true, token))
	}
}