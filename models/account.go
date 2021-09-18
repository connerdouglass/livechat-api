package models

import (
	"database/sql"
	"time"

	"github.com/connerdouglass/livechat-api/utils"
)

// Account is an admin account on the platform
type Account struct {
	ID           uint64 `gorm:"primaryKey"`
	Email        string
	PasswordSalt string
	PasswordHash string
	CreatedDate  time.Time
	DeletedDate  sql.NullTime
}

// VerifyPassword verifies a password on the account
func (a *Account) VerifyPassword(password string) bool {
	passwordHash := utils.Sha256Hex(a.PasswordSalt + password)
	return passwordHash == a.PasswordHash
}

// SetPassword sets a new password for the account
func (a *Account) SetPassword(password string) {
	a.PasswordSalt = utils.RandHexStrInt64()
	a.PasswordHash = utils.Sha256Hex(a.PasswordSalt + password)
}
