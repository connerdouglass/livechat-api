package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livechat-api/services"
)

func AppState(
	telegramService *services.TelegramService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Return the app state
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"telegram_bot_username": telegramService.BotUsername,
			},
		})

	}
}
