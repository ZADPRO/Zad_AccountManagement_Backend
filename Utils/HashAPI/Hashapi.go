package hashapi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Encrypt encrypts the given text using AES-256-CBC and PKCS7 padding.
// If encryptStatus is false, it returns the plain data wrapped in a map.
// The key is derived from ENCRYPT_API + token using SHA256.
func Encrypt(text interface{}, encryptStatus bool, token string) interface{} {
	if !encryptStatus {
		return text
	}

	// Derive 32-byte key from secret + token
	keyData := os.Getenv("ENCRYPT_API") + token
	key := sha256.Sum256([]byte(keyData)) // Always 32 bytes

	// Convert object to JSON string
	var plainText string
	switch v := text.(type) {
	case string:
		plainText = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to serialize object: %v", err)
		}
		plainText = string(bytes)
	}

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("failed to generate IV: %v", err)
	}

	// Create AES cipher and encrypt
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return fmt.Errorf("failed to create cipher: %v", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	padded := PKCS7Pad([]byte(plainText), aes.BlockSize)

	cipherText := make([]byte, len(padded))
	mode.CryptBlocks(cipherText, padded)

	// Return IV and encrypted hex string
	return []string{
		hex.EncodeToString(iv),
		hex.EncodeToString(cipherText),
	}
}
func Decrypt(encryptedData []string, token string) (interface{}, error) {
	if len(encryptedData) != 2 {
		return nil, errors.New("invalid encrypted data format")
	}

	// Derive the same key using ENCRYPT_API + token
	keyData := os.Getenv("ENCRYPT_API") + token
	key := sha256.Sum256([]byte(keyData)) // 32 bytes

	// Decode IV and cipherText from hex
	iv, err := hex.DecodeString(encryptedData[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %v", err)
	}
	cipherText, err := hex.DecodeString(encryptedData[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode cipherText: %v", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	// Check length
	if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("cipherText is not a multiple of the block size")
	}

	// Decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	plainPadded := make([]byte, len(cipherText))
	mode.CryptBlocks(plainPadded, cipherText)

	// Unpad using PKCS7
	plainText, err := PKCS7Unpad(plainPadded, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("unpadding error: %v", err)
	}

	// Try to decode JSON, otherwise return as string
	var result interface{}
	if err := json.Unmarshal(plainText, &result); err != nil {
		return string(plainText), nil
	}
	return result, nil
}

// PKCS7Unpad removes padding from the decrypted data
func PKCS7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padded data")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize {
		return nil, errors.New("invalid padding length")
	}
	for _, b := range data[len(data)-padLen:] {
		if b != byte(padLen) {
			return nil, errors.New("invalid padding byte")
		}
	}
	return data[:len(data)-padLen], nil
}

// PKCS7Pad adds padding to make data a multiple of block size
func PKCS7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padding := bytesRepeat(byte(padLen), padLen)
	return append(data, padding...)
}

func bytesRepeat(b byte, count int) []byte {
	result := make([]byte, count)
	for i := range result {
		result[i] = b
	}
	return result
}
