package Controller

import (
    "database/sql"
    "invoice-backend/Helper/HashAPI" 
    "invoice-backend/Services"
    "net/http"

    "github.com/gin-gonic/gin"
)

// GetDashboardStats aggregates high-level business metrics for the overview page.
// It retrieves summary data (e.g., total sales, pending invoices) and returns 
// an encrypted payload that can only be decrypted by the authenticated session holder.
func GetDashboardStats(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Session Context
		// Retrieve the JWT/Session token. This token acts as the 'Secret Key' 
		// for the AES/encryption algorithm during this specific request.
        token := getToken(c)

        // 2. Data Retrieval
		// Services layer handles the heavy lifting of SQL aggregation (SUM, COUNT, etc.).
        summary, err := Services.GetDashboardSummary(db)
        if err != nil {
            // SECURITY: Even error messages are encrypted. 
			// This prevents 'leaking' system info to attackers who don't have a valid token.
            c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
                "status":  false,
                "message": "Failed to load dashboard: " + err.Error(),
            }, true, token))
            return
        }

        // 3. Secure Transmission
		// The 'summary' object is converted to JSON and then encrypted via hashapi.
		// The React frontend interceptor will catch this, see the 'data' key, 
		// and use its locally stored token to turn it back into readable dashboard cards.
        c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{
            "status": true,
            "data":   summary,
        }, true, token))
    }
}