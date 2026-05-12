package Middleware

import (
	"strings"

	"invoice-backend/DB"
	accesstoken "invoice-backend/Helper/AccessToken"
	"invoice-backend/Helper/UserValidation"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract the Token from the Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		// Handle "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No "Bearer " prefix found
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		// 2. Validate the JWT (Check signature and expiration)
		token, err := accesstoken.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Invalid or expired token"})
			return
		}

		// 3. Extract Claims (Get the user ID stored in the token)
		claims, err := accesstoken.ExtractClaims(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Invalid claims"})
			return
		}

		userId := claims["user_id"]

		// 4. Live Status Check (Query the 'active_users' view via Helper)
		conn := DB.GetConn() // Get your GORM DB connection
		isActive, err := UserValidation.ValidateUser(conn, userId)
		if err != nil || !isActive {
			c.AbortWithStatusJSON(403, gin.H{"error": "Account is inactive or restricted"})
			return
		}

		// 5. Store data in context for the Controllers
		c.Set("user_id", userId)
		c.Set("user_role", claims["userType"])
		
		// ✅ CRITICAL: We store the raw token string.
		// The Controller needs this string to derive the key for HashAPI decryption/encryption.
		c.Set("raw_token", tokenString)

		c.Next()
	}
}