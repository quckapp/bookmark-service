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
	tagRepo := repository.NewTagRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	sharingRepo := repository.NewSharingRepository(db)
	noteRepo := repository.NewNoteRepository(db)
	reminderRepo := repository.NewReminderRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	previewRepo := repository.NewPreviewRepository(db)
	readLaterRepo := repository.NewReadLaterRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	versionRepo := repository.NewVersionRepository(db)
	expirationRepo := repository.NewExpirationRepository(db)
	templateRepo := repository.NewTemplateRepository(db)

	// Initialize services
	bookmarkService := service.NewBookmarkService(bookmarkRepo, folderRepo, redisClient, logger)
	folderService := service.NewFolderService(folderRepo, logger)
	tagService := service.NewTagService(tagRepo, logger)
	collectionService := service.NewCollectionService(collectionRepo, logger)
	sharingService := service.NewSharingService(sharingRepo, bookmarkRepo, logger)
	noteService := service.NewNoteService(noteRepo, logger)
	reminderService := service.NewBookmarkReminderService(reminderRepo, logger)
	favoriteService := service.NewFavoriteService(favoriteRepo, logger)
	analyticsService := service.NewBookmarkAnalyticsService(
		analyticsRepo, bookmarkRepo, folderRepo, tagRepo, collectionRepo, activityRepo, logger,
	)
	previewService := service.NewPreviewService(previewRepo, logger)
	readLaterService := service.NewReadLaterService(readLaterRepo, logger)
	commentService := service.NewCommentService(commentRepo, logger)
	versionService := service.NewVersionService(versionRepo, bookmarkRepo, logger)
	expirationService := service.NewExpirationService(expirationRepo, logger)
	templateService := service.NewTemplateService(templateRepo, bookmarkRepo, logger)

	// Initialize handlers
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkService)
	folderHandler := handler.NewFolderHandler(folderService)
	tagHandler := handler.NewTagHandler(tagService)
	collectionHandler := handler.NewCollectionHandler(collectionService)
	sharingHandler := handler.NewSharingHandler(sharingService)
	noteHandler := handler.NewNoteHandler(noteService)
	reminderHandler := handler.NewReminderHandler(reminderService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	previewHandler := handler.NewPreviewHandler(previewService)
	readLaterHandler := handler.NewReadLaterHandler(readLaterService)
	commentHandler := handler.NewCommentHandler(commentService)
	versionHandler := handler.NewVersionHandler(versionService)
	expirationHandler := handler.NewExpirationHandler(expirationService)
	templateHandler := handler.NewTemplateHandler(templateService)

	// Setup router with auth middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(goauth.RequestID())
	router.Use(goauth.CORS())

	// Health check (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "bookmark-service"})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	// JWT auth middleware for all /api/* routes
	authCfg := goauth.DefaultConfig(cfg.JWTSecret)
	authenticated := router.Group("")
	authenticated.Use(goauth.Auth(authCfg))

	// Bookmark CRUD
	api := authenticated.Group("/api/bookmarks")
	{
		api.POST("", bookmarkHandler.Create)
		api.GET("/:id", bookmarkHandler.GetByID)
		api.GET("/user/:userId", bookmarkHandler.GetByUser)
		api.GET("/user/:userId/workspace/:workspaceId", bookmarkHandler.GetByUserAndWorkspace)
		api.PUT("/:id", bookmarkHandler.Update)
		api.DELETE("/:id", bookmarkHandler.Delete)
		api.POST("/:id/move", bookmarkHandler.MoveToFolder)
		api.GET("/folder/:folderId", bookmarkHandler.GetByFolder)

		// Bulk operations
		api.POST("/bulk-delete", bookmarkHandler.BulkDelete)
		api.POST("/bulk-move", bookmarkHandler.BulkMove)
		api.POST("/reorder", bookmarkHandler.Reorder)

		// Tags on bookmarks
		api.POST("/:id/tags", tagHandler.TagBookmark)
		api.PUT("/:id/tags", tagHandler.ReplaceBookmarkTags)
		api.DELETE("/:id/tags/:tagId", tagHandler.UntagBookmark)
		api.GET("/:id/tags", tagHandler.GetBookmarkTags)

		// Notes on bookmarks
		api.POST("/:id/notes/:userId", noteHandler.Create)
		api.GET("/:id/notes", noteHandler.GetByBookmark)
		api.PUT("/:id/notes/:noteId", noteHandler.Update)
		api.DELETE("/:id/notes/:noteId", noteHandler.Delete)
		api.GET("/:id/notes/pinned", noteHandler.GetPinned)

		// Reminders on bookmarks
		api.POST("/:id/reminders/:userId", reminderHandler.Create)
		api.GET("/:id/reminders", reminderHandler.GetByBookmark)
		api.POST("/:id/reminders/:reminderId/cancel", reminderHandler.Cancel)
		api.DELETE("/:id/reminders/:reminderId", reminderHandler.Delete)

		// Favorites on bookmarks
		api.POST("/:id/favorite/:userId", favoriteHandler.Add)
		api.DELETE("/:id/favorite/:userId", favoriteHandler.Remove)
		api.GET("/:id/favorite/:userId", favoriteHandler.IsFavorite)

		// Activity on bookmarks
		api.GET("/:id/activity", analyticsHandler.GetBookmarkActivity)

		// Link Previews on bookmarks
		api.POST("/:id/preview", previewHandler.Generate)
		api.GET("/:id/preview", previewHandler.Get)

		// Comments on bookmarks
		api.POST("/:id/comments/:userId", commentHandler.Create)
		api.GET("/:id/comments", commentHandler.GetByBookmark)
		api.PUT("/:id/comments/:commentId", commentHandler.Update)
		api.DELETE("/:id/comments/:commentId", commentHandler.Delete)

		// Version History on bookmarks
		api.GET("/:id/versions", versionHandler.List)
		api.GET("/:id/versions/:versionId", versionHandler.Get)
		api.POST("/:id/versions/:versionId/restore", versionHandler.Restore)

		// Expiration on bookmarks
		api.POST("/:id/expiration/:userId", expirationHandler.Set)
		api.GET("/:id/expiration", expirationHandler.Get)
		api.DELETE("/:id/expiration", expirationHandler.Remove)
	}

	// Folders
	folders := authenticated.Group("/api/bookmark-folders")
	{
		folders.POST("", folderHandler.Create)
		folders.GET("/:id", folderHandler.GetByID)
		folders.GET("/user/:userId", folderHandler.GetByUser)
		folders.GET("/user/:userId/workspace/:workspaceId", folderHandler.GetByUserAndWorkspace)
		folders.PUT("/:id", folderHandler.Update)
		folders.DELETE("/:id", folderHandler.Delete)
		folders.POST("/reorder", folderHandler.Reorder)
	}

	// Tags
	tags := authenticated.Group("/api/bookmark-tags")
	{
		tags.POST("", tagHandler.Create)
		tags.GET("/:id", tagHandler.GetByID)
		tags.GET("/user/:userId", tagHandler.GetByUser)
		tags.GET("/user/:userId/workspace/:workspaceId", tagHandler.GetByUserAndWorkspace)
		tags.PUT("/:id", tagHandler.Update)
		tags.DELETE("/:id", tagHandler.Delete)
		tags.GET("/:id/bookmarks", tagHandler.GetBookmarksByTag)
	}

	// Collections
	collections := authenticated.Group("/api/bookmark-collections")
	{
		collections.POST("", collectionHandler.Create)
		collections.GET("/:id", collectionHandler.GetByID)
		collections.GET("/user/:userId", collectionHandler.GetByUser)
		collections.GET("/user/:userId/workspace/:workspaceId", collectionHandler.GetByUserAndWorkspace)
		collections.GET("/public/workspace/:workspaceId", collectionHandler.GetPublic)
		collections.PUT("/:id", collectionHandler.Update)
		collections.DELETE("/:id", collectionHandler.Delete)
		collections.POST("/:id/bookmarks", collectionHandler.AddBookmarks)
		collections.DELETE("/:id/bookmarks/:bookmarkId", collectionHandler.RemoveBookmark)
		collections.GET("/:id/bookmarks", collectionHandler.GetBookmarks)
	}

	// Sharing
	sharing := authenticated.Group("/api/bookmark-shares")
	{
		sharing.POST("/user/:userId", sharingHandler.Share)
		sharing.GET("/received/:userId", sharingHandler.GetSharedWithUser)
		sharing.GET("/sent/:userId", sharingHandler.GetSharedByUser)
		sharing.POST("/:id/accept", sharingHandler.Accept)
		sharing.POST("/:id/decline", sharingHandler.Decline)
		sharing.GET("/pending/:userId/count", sharingHandler.GetPendingCount)
	}

	// User-level resources
	users := authenticated.Group("/api/bookmark-users")
	{
		// Notes
		users.GET("/:userId/notes", noteHandler.GetByUser)

		// Reminders
		users.GET("/:userId/reminders", reminderHandler.GetByUser)
		users.GET("/:userId/reminders/pending", reminderHandler.GetPending)

		// Favorites
		users.GET("/:userId/favorites", favoriteHandler.GetFavorites)

		// Analytics and Stats
		users.GET("/:userId/stats", analyticsHandler.GetStats)
		users.GET("/:userId/recent", analyticsHandler.GetRecent)
		users.GET("/:userId/search", analyticsHandler.Search)
		users.GET("/:userId/duplicates", analyticsHandler.CheckDuplicate)
		users.GET("/:userId/export", analyticsHandler.Export)
		users.POST("/:userId/import/:workspaceId", analyticsHandler.Import)
		users.GET("/:userId/activity", analyticsHandler.GetActivity)
	}

	// Link Previews
	previews := authenticated.Group("/api/bookmark-previews")
	{
		previews.POST("", previewHandler.Create)
		previews.GET("/:url", previewHandler.GetByURL)
	}

	// Read Later
	readLater := authenticated.Group("/api/bookmark-readlater")
	{
		readLater.POST("/:userId", readLaterHandler.Add)
		readLater.GET("/:userId", readLaterHandler.List)
		readLater.PUT("/:id/status", readLaterHandler.UpdateStatus)
		readLater.GET("/:userId/stats", readLaterHandler.GetStats)
	}

	// Expirations
	expirations := authenticated.Group("/api/bookmark-expirations")
	{
		expirations.GET("/:userId/expiring", expirationHandler.GetExpiring)
	}

	// Templates
	templates := authenticated.Group("/api/bookmark-templates")
	{
		templates.POST("", templateHandler.Create)
		templates.GET("/user/:userId", templateHandler.GetByUser)
		templates.GET("/user/:userId/workspace/:workspaceId", templateHandler.GetByUserAndWorkspace)
		templates.PUT("/:id", templateHandler.Update)
		templates.DELETE("/:id", templateHandler.Delete)
		templates.POST("/:id/apply", templateHandler.Apply)
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
