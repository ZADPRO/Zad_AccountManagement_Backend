package Utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a one-way secure hash
func HashPassword(password string) (string, error) {
	// Cost 10 is the standard balance between speed and security
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash compares a raw password against the stored hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}