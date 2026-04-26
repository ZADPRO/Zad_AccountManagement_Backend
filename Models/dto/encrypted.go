package dto

// EncryptedPacket is the bridge for frontend-to-backend traffic.
// The frontend sends [ "iv_hex", "ciphertext_hex" ]
type EncryptedPacket struct {
	Data []string `json:"data" binding:"required"`
}