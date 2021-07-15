package services

import (
	"fmt"
	"strings"

	"github.com/godocompany/livechat-api/utils"
)

type TelegramUser struct {
	AuthDate  uint64  `json:"auth_date"`
	ID        uint64  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	PhotoUrl  *string `json:"photo_url"`
	Username  string  `json:"username"`
	Hash      string  `json:"hash"`
}

type TelegramService struct {
	BotAPIKey   string
	BotUsername string
}

// Verify verified the validity of a Telegram user
func (s *TelegramService) Verify(user *TelegramUser) bool {

	// Add all of the keyval pairs to a slice
	keyvals := []string{}
	keyvals = append(keyvals, fmt.Sprintf("auth_date=%d", user.AuthDate))
	if user.FirstName != nil {
		keyvals = append(keyvals, fmt.Sprintf("first_name=%s", *user.FirstName))
	}
	keyvals = append(keyvals, fmt.Sprintf("hash=%s", user.Hash))
	keyvals = append(keyvals, fmt.Sprintf("id=%d", user.ID))
	if user.LastName != nil {
		keyvals = append(keyvals, fmt.Sprintf("last_name=%s", *user.LastName))
	}
	if user.PhotoUrl != nil {
		keyvals = append(keyvals, fmt.Sprintf("photo_url=%s", *user.PhotoUrl))
	}
	keyvals = append(keyvals, fmt.Sprintf("username=%s", user.Username))

	// Create the data check string
	dataCheckString := strings.Join(keyvals, "\n")

	// Create the hash of the API key
	secretKey := utils.Sha256Hex(s.BotAPIKey)

	// Calculate the actual hash
	hash := utils.HmacSha256(dataCheckString, secretKey)

	// Compare the hashes
	return hash == user.Hash

}
