package v1

import (
	"github.com/connerdouglass/livechat-api/services"
	"github.com/connerdouglass/livechat-api/v1/hooks"
	"github.com/connerdouglass/livechat-api/v1/middleware"
	"github.com/gin-gonic/gin"
)

// Server is the API server instance
type Server struct {
	AccountsService   *services.AccountsService
	AuthTokensService *services.AuthTokensService
	ChatService       *services.ChatService
}

// Setup mounts the API server to the given group
func (s *Server) Setup(g *gin.RouterGroup) {

	// Register middleware for all routes
	g.Use(middleware.CheckAuth(s.AuthTokensService))

	// Register all of the public hooks that require no authentication
	s.setupPublicHooks(g)

	// Register authenticated hooks
	s.setupAuthenticatedHooks(g)

}

// setupPublicHooks mounts API hooks that are publicly accessible
func (s *Server) setupPublicHooks(g *gin.RouterGroup) {

	// Register public API routes
	g.POST("/app/get-state", hooks.AppState())
	g.POST("/auth/login", hooks.AuthLogin(
		s.AccountsService,
		s.AuthTokensService,
	))

}

// setupAuthenticatedHooks mounts API hooks that require account authentication
func (s *Server) setupAuthenticatedHooks(g *gin.RouterGroup) {

	// Require login for everything after this
	g.Use(middleware.RequireLogin())

	// Register authenticated API routes
	g.POST("/auth/whoami", hooks.AuthWhoAmI(
		s.AuthTokensService,
	))
	g.POST("/studio/chat/mute", hooks.StudioChatMute(
		s.AccountsService,
		s.ChatService,
	))
	g.POST("/studio/chat/unmute", hooks.StudioChatUnmute(
		s.AccountsService,
		s.ChatService,
	))

}
