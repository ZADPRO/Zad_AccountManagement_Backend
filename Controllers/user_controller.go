package Controller

import (
	"database/sql"
	"encoding/json"
	hashapi "invoice-backend/Helper/HashAPI"
	"invoice-backend/Models/dto"
	"invoice-backend/Models/responses"
	"invoice-backend/Services"
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
)


// CreateUser handles the registration of a new user in the system.
// It decrypts the request, passes it to the Service layer (which handles 
// hashing and email notifications), and returns an encrypted response.
func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		var packet dto.EncryptedPacket
		
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Packet"}, true, token))
			return
		}

		// 1. Decrypt incoming data using the session token
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Decryption failed"}, true, token))
			return
		}

		var req dto.CreateUserRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		// 2. Logic: The service layer adds the user to the DB and triggers a Welcome Email
		id, email, tempPassword, err := Services.AddNewUser(db, req)
		if err != nil {
		// Handle specific business logic error: Email duplication
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, hashapi.Encrypt(gin.H{
				"status": false,
				"message": "email already exists",
			}, true, token))
			return
    }

    c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
        "status": false,
        "message": err.Error(),
    }, true, token))
    return
}

fmt.Println("EMAIL:", email)
fmt.Println("TEMP PASSWORD:", tempPassword)

		// 3. Return success with the new UserID
		resp := hashapi.Encrypt(gin.H{
	"status": true,
	"message": "User created successfully",
	"userId": id,
	"email": email,
	"tempPassword": tempPassword,
}, true, token)
		c.JSON(http.StatusCreated, resp)
	}
}

// GetUserList retrieves all user records for administrative display.
func GetUserList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		data, err := Services.FetchAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(responses.UserListResponse{
			BaseResponse: dto.BaseResponse{Status: true, Message: "Users fetched successfully"},
			Users:        data,
		}, true, token))
	}
}

// GetUserByID fetches a specific user's details based on the URL parameter ID.
func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		idStr := c.Param("id")
		userID, _ := strconv.Atoi(idStr)

		user, err := Services.GetUserByID(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, hashapi.Encrypt(gin.H{"status": false, "message": "User not found"}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "data": user}, true, token))
	}
}

// GetProfile retrieves the profile details of the currently authenticated user.
func GetProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		// Standardize this key with your AuthMiddleware (usually "user_id")
		val, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Unauthorized"}, true, token))
			return
		}
		adminID := int(val.(float64))
		profile, err := Services.GetUserProfile(db, adminID)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Profile Error"}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{"status": true, "data": profile}, true, token))
	}
}

// UpdateUser modifies profile or user information.
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		idStr := c.Param("id")
		userID, _ := strconv.Atoi(idStr)

		var packet dto.EncryptedPacket
		c.ShouldBindJSON(&packet)

		decrypted, _ := hashapi.Decrypt(packet.Data, token)
		var req dto.UpdateUserRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		err := Services.UpdateUser(db, userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Update failed"}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.BaseResponse{Status: true, Message: "User updated successfully"}, true, token))
	}
}

// ForceChangePassword allows an admin or system to reset a specific user's password.
// This is typically used for "Forgot Password" or initial setup scenarios.
func ForceChangePassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		var packet dto.EncryptedPacket
		c.ShouldBindJSON(&packet)

		decrypted, _ := hashapi.Decrypt(packet.Data, token)

		// Helper struct for decryption mapping
		var req struct {
			NewPassword string `json:"password"`
			UserID      int    `json:"userId"`
		}
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		err := Services.UpdateUserPassword(db, req.UserID, req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.BaseResponse{Status: true, Message: "Password updated!"}, true, token))
	}
}

// DeleteUser performs a deletion of a user record.
// Logic Note: The adminID is passed to the service to audit who performed the deletion.
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		adminIDRaw, _ := c.Get("user_id")
		adminID := int(adminIDRaw.(float64))

		idStr := c.Param("id")
		userID, _ := strconv.Atoi(idStr)

		err := Services.DeleteUser(db, userID, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Delete failed"}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.BaseResponse{Status: true, Message: "User deleted"}, true, token))
	}
}