package utils

import (
	"github.com/connerdouglass/livechat-api/models"
	"github.com/gin-gonic/gin"
)

// CtxGetAccount gets the account (or nil) from a Gin context
func CtxGetAccount(c *gin.Context) *models.Account {

	// Get the account from the context
	owner, exists := c.Get("account")
	if !exists || owner == nil {
		return nil
	}

	// Perform a typecheck on the account
	account, ok := owner.(*models.Account)
	if !ok || account == nil {
		return nil
	}

	// Return the account
	return account

}
