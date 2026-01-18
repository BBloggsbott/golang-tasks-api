package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BBloggsbott/task-api/internal/config"
	"github.com/BBloggsbott/task-api/internal/database"
	"github.com/BBloggsbott/task-api/internal/handlers"
	"github.com/BBloggsbott/task-api/internal/repository"
	"github.com/BBloggsbott/task-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// DEBUG: Print config (hide password)
	log.Printf("Config - Host: %s, Port: %s, User: %s, Password length: %d, Database: %s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		len(cfg.Database.Password), // Show length, not actual password
		cfg.Database.Database,
	)

	// Connect to database
	db, err := database.NewMySQLDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	log.Println("Successfully connected to database")

	// Initialize layers
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	// Setup router
	router := handlers.SetupRouter(taskService)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
