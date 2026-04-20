package Controller

import (
	"database/sql"
	"invoice-backend/Model"
	"invoice-backend/Service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Model.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Invalid Data: " + err.Error()})
			return
		}

		// Service handles: Temp Password -> Hashing -> Encryption -> DB Insert
		id, err := Service.AddNewUser(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Model.CreateUserResponse{
			BaseResponse: Model.BaseResponse{Status: true, Message: "User created successfully"},
			UserID:       id,
		})
	}
}

func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, _ := strconv.Atoi(idStr)

		var req Model.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Invalid Request"})
			return
		}

		// Service handles: Encryption -> DB Update
		err := Service.UpdateUser(db, userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Update failed"})
			return
		}

		c.JSON(http.StatusOK, Model.BaseResponse{Status: true, Message: "User updated successfully"})
	}
}

func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, _ := strconv.Atoi(idStr)

		// Service handles: DB Fetch -> Decryption
		user, err := Service.GetUserByID(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, Model.BaseResponse{Status: false, Message: "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   user, // Decrypted data sent to frontend
		})
	}
}

func GetProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ⚠️ IMPORTANT: Verify if your middleware uses "userId" or "userID"
		val, exists := c.Get("userId") 
		if !exists {
			c.JSON(http.StatusUnauthorized, Model.BaseResponse{Status: false, Message: "Unauthorized"})
			return
		}

		// Service handles: DB Fetch -> Decryption
		profile, err := Service.GetUserProfile(db, val.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Profile Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   profile,
		})
	}
}

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, _ := strconv.Atoi(idStr)

		var req Model.DeleteUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Admin ID required"})
			return
		}

		err := Service.DeleteUser(db, userID, req.AdminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Delete failed"})
			return
		}

		c.JSON(http.StatusOK, Model.BaseResponse{Status: true, Message: "User deleted"})
	}
}

func ForceChangePassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NewPassword string `json:"password" binding:"required,min=8"`
			UserID      int    `json:"userId" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{
				Status:  false,
				Message: "Password must be at least 8 characters",
			})
			return
		}

		// Call the service to update the password and flip the flag
		err := Service.UpdateUserPassword(db, req.UserID, req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{
				Status:  false,
				Message: "Failed to update password: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Model.BaseResponse{
			Status:  true,
			Message: "Password updated successfully! Please login again.",
		})
	}
}
