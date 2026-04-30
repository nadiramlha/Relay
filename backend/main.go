package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/nadiramlha/relay/backend/agent"
	"github.com/nadiramlha/relay/backend/db"
	"github.com/nadiramlha/relay/backend/handlers"
	"github.com/nadiramlha/relay/backend/middleware"
)

func main() {
	// Load .env file (non-fatal if absent in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables from system")
	}

	// Connect to PostgreSQL
	ctx := context.Background()
	database, err := db.New(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Connected to PostgreSQL")

	// Initialise AI agent
	aiAgent, err := agent.New()
	if err != nil {
		log.Fatalf("Failed to initialise AI agent: %v", err)
	}
	log.Println("AI agent initialised")

	// Build router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	h := handlers.New(database, aiAgent)

	// Routes
	api := r.Group("/api")
	{
		api.GET("/health", h.HealthCheck)
		api.POST("/conversations", h.CreateConversation)
		api.GET("/conversations", h.GetConversations)
		api.DELETE("/conversations/:id", h.DeleteConversation)
		api.GET("/conversations/:id/messages", h.GetMessages)
		api.POST("/chat", h.Chat)
	}

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
