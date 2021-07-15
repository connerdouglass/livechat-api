package utils

import (
	"errors"
	"strings"
)

// ValidatePassword checks if a password is valid to be used for an account
func ValidatePassword(password string) error {

	// If the password starts or ends with whitespace
	if strings.TrimSpace(password) != password {
		return errors.New("password cannot begin or end with whitespace")
	}

	// If the length of the password is too short
	if len(password) < 3 {
		return errors.New("password must be at least 3 characters long")
	}

	// If the password is too long
	if len(password) > 64 {
		return errors.New("password must not be longer than 64 characters")
	}

	// Return no error otherwise
	return nil

}
