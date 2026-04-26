package Utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

// getDBKey derives a 32-byte key from your static DB secret
func getDBKey() []byte {
	hash := sha256.Sum256([]byte(os.Getenv("DB_ENCRYPT_SECRET")))
	return hash[:]
}

// EncryptForDB uses Deterministic GCM (Same input = Same ciphertext)
func EncryptForDB(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	block, err := aes.NewCipher(getDBKey())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 🔒 DETERMINISTIC STEP: Nonce is derived from the text itself
	// This makes the field searchable in the database
	h := sha256.New()
	h.Write([]byte(plainText))
	nonce := h.Sum(nil)[:gcm.NonceSize()]

	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)
	
	// Prepend nonce so Decrypt can read it
	return base64.StdEncoding.EncodeToString(append(nonce, cipherText...)), nil
}

// DecryptFromDB reverses the storage encryption
func DecryptFromDB(cryptoText string) (string, error) {
	if cryptoText == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(getDBKey())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}