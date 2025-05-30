package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ikurotime/code-engine/internal/handlers"
	"ikurotime/code-engine/internal/services"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize services
	executor := services.NewExecutor(1, 30*time.Second, logger)
	handler := handlers.NewHandler(executor, logger)

	// Setup routes
	router := http.NewServeMux()
	router.HandleFunc("/health", handler.HealthCheck)
	router.HandleFunc("/execute", handler.Execute)

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	logger.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Server failed to start:", err)
	}
}
