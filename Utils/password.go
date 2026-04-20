package Utils

import (
    "crypto/rand"
    "encoding/base64"
)

// GenerateRandomPassword creates a secure 12-character string
func GenerateRandomPassword() string {
    b := make([]byte, 9)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
} 

