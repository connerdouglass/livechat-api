package hooks

import (
	"net/http"

	"github.com/connerdouglass/livechat-api/services"
	"github.com/gin-gonic/gin"
)

type StudioChatMuteReq struct {
	OrganizationID uint64                `json:"organization_id"`
	User           services.ChatUserInfo `json:"user"`
}

func StudioChatMute(
	accountsService *services.AccountsService,
	chatService *services.ChatService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioChatMuteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the account sending the request
		// account := utils.CtxGetAccount(c)

		// Mute the user on the chat
		if _, err := chatService.MuteUser(req.OrganizationID, &req.User, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Otherwise return something successfully
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
