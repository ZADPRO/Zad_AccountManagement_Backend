package Controller

import (
	"database/sql"
	"encoding/json"
	"invoice-backend/Helper/HashAPI"
	"invoice-backend/Models/dto"
	"invoice-backend/Models/responses"
	"invoice-backend/Services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		var packet dto.EncryptedPacket

		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Packet"}, true, token))
			return
		}

		// 1. Decrypt Axios Request
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Decryption failed"}, true, token))
			return
		}

		var req dto.CreateUserRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		// 2. Service handles DB Encryption & Welcome Email
		id, err := Services.AddNewUser(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}

		// 3. Encrypt Response for Frontend
		resp := hashapi.Encrypt(dto.CreateUserResponse{
			BaseResponse: dto.BaseResponse{Status: true, Message: "User created successfully"},
			UserID:       id,
		}, true, token)
		c.JSON(http.StatusCreated, resp)
	}
}

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