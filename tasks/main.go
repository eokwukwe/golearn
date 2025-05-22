package main

import (
	"log"
	"net/http"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/handlers"
)

func main() {
	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Define routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	http.HandleFunc("/api/v1/register", handlers.Register)
	http.HandleFunc("/api/v1/login", handlers.Login)

	// Start server
	log.Printf("Starting server on :7070")
	if err := http.ListenAndServe(":7070", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
