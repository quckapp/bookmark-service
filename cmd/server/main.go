package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	goauth "github.com/quckapp/go-auth"
	"github.com/quckapp/bookmark-service/internal/config"
	"github.com/quckapp/bookmark-service/internal/handler"
	"github.com/quckapp/bookmark-service/internal/repository"
	"github.com/quckapp/bookmark-service/internal/service"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient := config.InitRedis(cfg)

	// Initialize repositories
	bookmarkRepo := repository.NewBookmarkRepository(db)
	folderRepo := repository.NewFolderRepository(db)

	// Initialize services
	bookmarkService := service.NewBookmarkService(bookmarkRepo, folderRepo, redisClient, logger)
	folderService := service.NewFolderService(folderRepo, logger)

	// Initialize handlers
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	folderHandler := handler.NewFolderHandler(folderService)

	// Setup router with shared middleware
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(goauth.RequestID())
	router.Use(goauth.CORS())
	router.Use(goauth.Logger(logrusLogger))

	// Health check (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	// Auth middleware configuration
	authCfg := goauth.DefaultConfig(cfg.JWTSecret)

	// Authenticated API routes
	authenticated := router.Group("")
	authenticated.Use(goauth.Auth(authCfg))

	api := authenticated.Group("/api/bookmarks")
	{
		api.POST("", bookmarkHandler.Create)
		api.GET("/:id", bookmarkHandler.GetByID)
		api.GET("/user/:userId", bookmarkHandler.GetByUser)
		api.GET("/user/:userId/workspace/:workspaceId", bookmarkHandler.GetByUserAndWorkspace)
		api.PUT("/:id", bookmarkHandler.Update)
		api.DELETE("/:id", bookmarkHandler.Delete)
		api.POST("/:id/move", bookmarkHandler.MoveToFolder)
	}

	folders := authenticated.Group("/api/bookmark-folders")
	{
		folders.POST("", folderHandler.Create)
		folders.GET("/:id", folderHandler.GetByID)
		folders.GET("/user/:userId", folderHandler.GetByUser)
		folders.PUT("/:id", folderHandler.Update)
		folders.DELETE("/:id", folderHandler.Delete)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5010"
	}

	logger.Info("Starting bookmark service", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
