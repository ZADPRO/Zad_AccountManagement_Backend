package hashapi

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os" 
	"errors"
)

// Encrypt prepares transient data for the frontend.
// It uses AES-256-CBC with PKCS7 padding.
func Encrypt(data interface{}, encryptStatus bool, token string) interface{} {
	if !encryptStatus {
		return data
	}

	// 1. Derive a unique 32-byte key from the API Secret and the JWT Token
	keyData := os.Getenv("ENCRYPT_API") + token 

	key := sha256.Sum256([]byte(keyData))
	
	// 2. Serialize the data to JSON
	var plainText []byte
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("Serialization error: %v", err)
	}
	plainText = bytes

	// 3. Create the AES block cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return fmt.Sprintf("Cipher error: %v", err)
	}

	// 4. Generate a random Initialization Vector (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Sprintf("IV generation error: %v", err)
	}

	// 5. Apply PKCS7 Padding
	paddedText := pkcs7Pad(plainText, aes.BlockSize)

	// 6. Encrypt using CBC mode
	cipherText := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddedText)

	// 7. Return as [iv_hex, ciphertext_hex] for the frontend
	return []string{
		hex.EncodeToString(iv),
		hex.EncodeToString(cipherText),
	}
}

// pkcs7Pad adds standard PKCS7 padding to the plaintext
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
} 


// Decrypt handles data coming FROM the frontend.
// encryptedData is the [iv_hex, ciphertext_hex] array.
func Decrypt(encryptedData []string, token string) (interface{}, error) {
	if len(encryptedData) != 2 {
		return nil, errors.New("invalid encrypted data format")
	}

	// 1. Reconstruct the same key used by the frontend: Secret + Token
	keyData := os.Getenv("ENCRYPT_API") + token
	key := sha256.Sum256([]byte(keyData))

	// 2. Decode hex values
	iv, err := hex.DecodeString(encryptedData[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %v", err)
	}
	cipherText, err := hex.DecodeString(encryptedData[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	// 3. Initialize AES-CBC
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	// 4. Decrypt blocks
	mode := cipher.NewCBCDecrypter(block, iv)
	plainPadded := make([]byte, len(cipherText))
	mode.CryptBlocks(plainPadded, cipherText)

	// 5. Remove PKCS7 Padding
	plainText, err := pkcs7Unpad(plainPadded, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("unpadding error: %v", err)
	}

	// 6. Try to unmarshal JSON, otherwise return as string
	var result interface{}
	if err := json.Unmarshal(plainText, &result); err != nil {
		return string(plainText), nil
	}
	return result, nil
}

// pkcs7Unpad removes padding bytes from decrypted data
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padded data length")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize {
		return nil, errors.New("invalid padding length")
	}
	// Verify padding bytes
	for _, b := range data[len(data)-padLen:] {
		if b != byte(padLen) {
			return nil, errors.New("invalid padding byte values")
		}
	}
	return data[:len(data)-padLen], nil
}