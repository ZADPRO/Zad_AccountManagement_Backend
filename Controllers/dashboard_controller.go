package Controller

import (
    "database/sql"
    "invoice-backend/Helper/HashAPI" // Your encryption helper
    "invoice-backend/Services"
    "net/http"

    "github.com/gin-gonic/gin"
)

func GetDashboardStats(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get the session token for encryption
        token := getToken(c)

        // 2. Fetch stats from service
        summary, err := Services.GetDashboardSummary(db)
        if err != nil {
            // Encrypt error response so Axios doesn't crash
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
                "status":  false,
                "message": "Failed to load dashboard: " + err.Error(),
            }, true, token))
            return
        }

        // 3. Encrypt the data payload
        // This allows your React component to decrypt it into the dashboard cards
        c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{
            "status": true,
            "data":   summary,
        }, true, token))
    }
}