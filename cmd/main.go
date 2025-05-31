package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ikurotime/code-engine/config"
	"ikurotime/code-engine/internal/handlers"
	"ikurotime/code-engine/internal/services"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize services
	executor := services.NewExecutor(config.Server.MaxConcurrentExecutions, time.Duration(config.Server.ExecutionTimeout)*time.Second, logger)
	handler := handlers.NewHandler(executor, logger)

	// Setup routes
	router := http.NewServeMux()
	router.HandleFunc("/", handler.Home)
	router.HandleFunc("/health", handler.HealthCheck)
	router.HandleFunc("/execute", handler.Execute)

	// Create server
	server := &http.Server{
		Addr:    config.Server.Port,
		Handler: router,
	}

	// Channel to listen for interrupt signal to terminate gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Println("Shutdown signal received")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	logger.Println("Shutting down HTTP server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("Server forced to shutdown: %v", err)
	} else {
		logger.Println("HTTP server shutdown completed")
	}

	// Shutdown executor and cleanup containers
	logger.Println("Shutting down executor and cleaning up containers...")
	executor.Shutdown()

	logger.Println("Graceful shutdown completed")
}
