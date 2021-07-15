package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hex calculates the SHA-256 hash of the input string, and returns the hash as a hexadecimal-encoded string
func Sha256Hex(input string) string {

	// Calculate the SHA256 hash, which returns an array of bytes
	hashBytes := sha256.Sum256([]byte(input))

	// Encode it as a hex string and return
	return hex.EncodeToString(hashBytes[:])

}

func HmacSha256(input, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	sum := h.Sum([]byte(input))
	return hex.EncodeToString(sum)
}
