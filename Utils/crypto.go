package Utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"encoding/hex" 
	"crypto/sha256"
)

func getEncryptionKey() ([]byte, error) {
    keyStr := os.Getenv("ENCRYPT_API")
    if keyStr == "" {
        return nil, fmt.Errorf("ENCRYPT_API environment variable is not set")
    }

    // Decode the 64-character hex string into 32 raw bytes
    key, err := hex.DecodeString(keyStr)
    if err != nil {
        return nil, fmt.Errorf("failed to decode hex key: %v", err)
    }

    if len(key) != 32 {
        return nil, fmt.Errorf("decoded ENCRYPT_API must be 32 bytes (current: %d)", len(key))
    }

    return key, nil
}

func Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cryptoText string) (string, error) {
	if cryptoText == "" {
		return "", nil
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
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
		return "", fmt.Errorf("decrypt error (wrong key?): %v", err)
	}

	return string(plainText), nil
} 


func EncryptDeterministic(plainText string) (string, error) {
    if plainText == "" { return "", nil }

    key, err := getEncryptionKey()
    if err != nil { return "", err }

    block, err := aes.NewCipher(key)
    if err != nil { return "", err }

    gcm, err := cipher.NewGCM(block)
    if err != nil { return "", err }

    // 🔒 DETERMINISTIC STEP:
    // We derive the nonce from the plaintext itself.
    // Same input + Same key = Same Nonce = Same Ciphertext.
    h := sha256.New()
    h.Write([]byte(plainText))
    nonce := h.Sum(nil)[:gcm.NonceSize()] 

    // Seal the data
    cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)
    
    // Prepend nonce so the existing Decrypt function can still read it!
    return base64.StdEncoding.EncodeToString(append(nonce, cipherText...)), nil
}