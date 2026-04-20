package Controller

import (
    "database/sql"
    "invoice-backend/Model"
    "invoice-backend/Utils"
    "net/http"
    "github.com/gin-gonic/gin" 
    "invoice-backend/Service" 
	"fmt"
    
)

func Login(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req Model.LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Invalid payload"})
            return
        }

        userID, roleName, userName, isFirstLogin, err := Service.VerifyLogin(db, req.Email, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, Model.BaseResponse{Status: false, Message: err.Error()})
            return
        }

        // 🔥 This generates the NEW token.
        token, err := Utils.GenerateToken(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Token error"})
            return
        }

        // 🟢 PREPARE DATA
        loginData := Model.LoginResponse{
            BaseResponse: Model.BaseResponse{
                Status:  true,
                Message: "Login successful",
            },
            Token:        token,
            Role:         roleName,
            Username:     userName,
            IsFirstLogin: isFirstLogin,
            UserId:       userID,
        }

        // 🔐 ENCRYPT (Use the fresh token as the key)
        // If you want the login response encrypted:
        // encryptedBody := hashapi.Encrypt(loginData, true, token)
        
        // For now, let's keep Login plain so React can easily read the token
        // but ensure we return the fresh token clearly.
        c.JSON(http.StatusOK, loginData)
    }
}

func ChangePassword(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {

        var req struct {
            UserID      int    `json:"userid" binding:"required"`
            NewPassword string `json:"newPassword" binding:"required"`
        }

        // 🔴 prevents your previous 400 error
        if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("Binding Error: %v\n", err)
            c.JSON(http.StatusBadRequest, gin.H{
                "status": false,
                "message": "Invalid request payload",
            })
            return
        }

        err := Service.ChangePassword(db, req.UserID, req.NewPassword)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": false,
                "message": err.Error(),
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "status":  true,
            "message": "Password updated successfully",
        })
    }
}