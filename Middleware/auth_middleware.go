package Middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("❌ Rejected: No Authorization Header")
			c.AbortWithStatusJSON(401, gin.H{"message": "Missing token"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("❌ Rejected: Invalid Header Format")
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token format"})
			return
		}

		tokenString := parts[1]
		secretKey := []byte(os.Getenv("JWT_SECRET"))

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// 🔍 DETAILED ERROR CHECKING
		if err != nil {
			fmt.Printf("❌ JWT Error: %v\n", err) // This prints the EXACT error to your console
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized: " + err.Error()})
			return
		}

		if !token.Valid {
			fmt.Println("❌ Rejected: Token is not valid")
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
			return
		}

		// ✅ CRITICAL FOR ENCRYPTION: 
		// You MUST set the token here so the Controller can use it as a key
		c.Set("token", tokenString)

		fmt.Println("✅ Middleware: Token validated and passed to context")
		c.Next()
	}
}