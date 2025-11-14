package main

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/lang-portal/backend/internal/database"
	"github.com/lang-portal/backend/internal/handlers"
	"github.com/lang-portal/backend/internal/middleware"
)

func main() {
	// Initialize database
	dbPath := "words.db"
	err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	log.Println("Database connected successfully")

	// Setup Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.SetupCORS())
	router.Use(middleware.Logger())

	// API routes
	api := router.Group("/api")
	{
		// Dashboard endpoints
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/last_study_session", handlers.GetLastStudySession)
			dashboard.GET("/study_progress", handlers.GetStudyProgress)
			dashboard.GET("/quick_stats", handlers.GetQuickStats)
		}

		// Words endpoints
		api.GET("/words", handlers.GetWords)
		api.GET("/words/:id", handlers.GetWordByID)

		// Groups endpoints
		api.GET("/groups", handlers.GetGroups)
		api.GET("/groups/:id", handlers.GetGroupByID)
		api.GET("/groups/:id/words", handlers.GetGroupWords)
		api.GET("/groups/:id/study_sessions", handlers.GetGroupStudySessions)

		// Study sessions endpoints
		api.GET("/study_sessions", handlers.GetStudySessions)
		api.GET("/study_sessions/:id", handlers.GetStudySessionByID)
		api.GET("/study_sessions/:id/words", handlers.GetStudySessionWords)
		api.POST("/study_sessions/:id/words/:word_id/review", handlers.RecordWordReview)

		// Study activities endpoints
		api.GET("/study_activities/:id", handlers.GetStudyActivityByID)
		api.GET("/study_activities/:id/study_sessions", handlers.GetActivityStudySessions)
		api.POST("/study_activities", handlers.CreateStudySession)

		// Utility endpoints
		api.POST("/reset_history", handlers.ResetHistory)
		api.POST("/full_reset", func(c *gin.Context) {
			migrationsDir := filepath.Join("db", "migrations")
			err := database.FullReset(migrationsDir)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"success": true, "message": "Full database reset completed successfully"})
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Language Learning Portal API"})
	})

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
