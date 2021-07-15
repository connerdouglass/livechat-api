package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livechat-api/services"
)

// CheckAuth creates a middleware function that parses auth token header and adds the account to the context
func CheckAuth(authTokensService *services.AuthTokensService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initially, store nil in the context
		c.Set("bearer_token", nil)
		c.Set("account", nil)

		// Get the authorization header, trimmed
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		// If it's empty skip
		if len(authHeader) == 0 {
			c.Next()
			return
		}

		// If it doesn't have the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		// Parse out the bearer token
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		c.Set("bearer_token", token)

		// Find the account of the token
		account, err := authTokensService.GetAccountForToken(token)
		if err != nil {
			// fmt.Println("auth token error: ", err)
			c.Next()
			return
		}
		if account != nil {
			c.Set("account", account)
		}

		// Move to the next
		c.Next()

	}
}
