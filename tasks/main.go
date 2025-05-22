package main

import (
	"log"
	"net/http"

	"github.com/eokwukwe/golearn/tasks/config"
	"github.com/eokwukwe/golearn/tasks/handlers"
	"github.com/eokwukwe/golearn/tasks/middleware"
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

	// Handle GET and POST requests for tasks separately
	http.HandleFunc("/api/v1/tasks", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetTasks(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreateTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Handle GET request for a single task
	http.HandleFunc("/api/v1/tasks/{id}", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetOneTask(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.DeleteTask(w, r)
		} else if r.Method == http.MethodPut {
			handlers.UpdateTask(w, r)
		} else if r.Method == http.MethodPatch {
			handlers.CompleteTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Start server
	log.Printf("Starting server on :7070")
	if err := http.ListenAndServe(":7070", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
