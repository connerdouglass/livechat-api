package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/godocompany/livechat-api/models"
	"github.com/godocompany/livechat-api/services"
	v1 "github.com/godocompany/livechat-api/v1"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	//================================================================================
	// Create the database connection
	//================================================================================

	// Get the datbase driver for the database string
	dbDriver := ParseDatabaseDriver(os.Getenv("DB_URL"))
	if dbDriver == nil {
		log.Fatalln("Failed to create database driver. Check DB_URL environment variable")
	}

	// Create the database connection
	db, err := gorm.Open(dbDriver, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.Account{},
		&models.Badge{},
		&models.ChatRoom{},
		&models.MutedUser{},
		&models.Organization{},
	)

	//================================================================================
	// Setup the WebSockets server
	//================================================================================

	// Get all of the allowed origins
	allowedOrigins := GetAllowedOrigins()

	// Create the server
	socketIoServer := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: checkOrigin(allowedOrigins),
			},
			&websocket.Transport{
				CheckOrigin: checkOrigin(allowedOrigins),
			},
		},
	})
	go socketIoServer.Serve()

	//================================================================================
	// Create all the service instances
	//================================================================================

	// Create the rest of the services
	telegramService := &services.TelegramService{
		BotAPIKey:   os.Getenv("TELEGRAM_BOT_API_KEY"),
		BotUsername: os.Getenv("TELEGRAM_BOT_USERNAME"),
	}
	chatService := &services.ChatService{
		DB: db,
	}
	socketsService := &services.SocketsService{
		Server:          socketIoServer,
		TelegramService: telegramService,
		ChatService:     chatService,
	}
	accountsService := &services.AccountsService{DB: db}
	authTokensService := &services.AuthTokensService{
		DB:            db,
		SigningPepper: os.Getenv("AUTH_TOKEN_SIGNING_PEPPER"),
	}

	// Do some final update on the sockets service
	// Needed because it has a circular relationship with other services
	socketsService.Setup()

	//================================================================================
	// Setup the Gin HTTP router
	//================================================================================

	// Create the Gin router
	r := gin.Default()

	// Configure CORS for the API
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = GetAllowedOrigins()
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Accept", "User-Agent", "Authorization")
	r.Use(cors.New(corsCfg))

	// Create the API instance
	api := &v1.Server{
		AccountsService:   accountsService,
		AuthTokensService: authTokensService,
		TelegramService:   telegramService,
		ChatService:       chatService,
	}

	// Mount the API routes
	api.Setup(r.Group("v1"))

	// Create a mux to serve both the HTTP and Socket.IO servers
	mux := http.NewServeMux()
	mux.Handle("/socket.io/", socketIoServer)
	mux.Handle("/", r)

	// Run the server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicln(err)
	}

}

// GetAllowedOrigins gets the slice of allowed CORS origins
func GetAllowedOrigins() []string {

	// Get the list of origins allowed
	env, ok := os.LookupEnv("CORS_ALLOW_ORIGINS")
	if !ok {
		return []string{}
	}

	// Create the slice for it
	origins := []string{}

	// Split up the env value
	originsRaw := strings.Split(env, ",")
	for _, originRaw := range originsRaw {
		origin := strings.TrimSpace(originRaw)
		origins = append(origins, origin)
	}

	// Return the origins slice
	return origins

}
