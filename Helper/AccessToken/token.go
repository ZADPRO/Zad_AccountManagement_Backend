package accesstoken

// Package accesstoken provides utilities for generating, parsing, and validating JWTs.
// It relies on the "JWT_SECRET" environment variable for cryptographic signing.

import (
	"fmt"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateToken generates a JWT for a user. 
// The token string is used later as a dynamic salt for HashAPI encryption.
func CreateToken(userId interface{}) (string, error) {
	// Retrieve secret from environment for security. 
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	
	// Define the payload (Claims). 
	// Standard 'exp' claim ensures the token is automatically invalid after 24 hours.
	claims := jwt.MapClaims{
		"user_id":  userId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiration
	}

	// Create the token using the HS256 algorithm (Symmetric signing).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
} 


// getAdminID is a robust helper to extract the user ID from the Gin context.
// Because JSON and JWT libraries often decode numbers as float64, this function 
// handles multiple numeric types to prevent type assertion panics.
func getAdminID(c *gin.Context) int {
    val, exists := c.Get("user_id")
    if !exists {
        return 0
    } 
	// Type switch to safely handle numeric variations coming from the token claims.
    switch v := val.(type) {
    case float64:
        return int(v)
    case int:
        return v
    case int64:
        return int(v)
    default:
        return 0
    }
} 

// ValidateJWT checks the token signature and expiration
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
} 

// ExtractClaims unpackages the data stored inside a validated token.
// Use this to get the "user_id" or "exp" values after calling ValidateJWT.
func ExtractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}