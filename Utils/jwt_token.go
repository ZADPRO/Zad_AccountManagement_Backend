package Utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT for a specific UserID
func GenerateToken(userID int) (string, error) {
	// 1. Get the secret key from environment variables
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// 2. Create the Claims (the payload of the token)
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), 
	}

	// 3. Create the token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Sign the token with our secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}