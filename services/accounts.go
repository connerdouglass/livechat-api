package services

import (
	"errors"

	"github.com/godocompany/livechat-api/models"
	"gorm.io/gorm"
)

// AccountsService manages account access to the dashboard. This is admin access, and is unrelated to
// Telegram accounts in the live chat itself
type AccountsService struct {
	DB *gorm.DB
}

// GetAccountByEmail gets the account with the provided email address
func (s *AccountsService) GetAccountByEmail(email string) (*models.Account, error) {
	var account models.Account
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("email LIKE ?", email).
		First(&account).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// FindByLogin finds an account with the provided login credentials
func (s *AccountsService) FindByLogin(email, password string) (*models.Account, error) {

	// Find the account with the email
	var account models.Account
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("email LIKE ?", email).
		First(&account).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Verify the password
	if !account.VerifyPassword(password) {
		return nil, nil
	}

	// Return the account
	return &account, nil

}
