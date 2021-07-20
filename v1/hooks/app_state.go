package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppState() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Return the app state
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
