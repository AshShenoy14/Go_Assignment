package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize database
	db, err := NewDatabase("users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handler
	handler := NewHandler(db)

	// Setup routes
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/users", handler.GetUsersHandler)
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetUserHandler(w, r)
		} else if r.Method == http.MethodPut {
			handler.UpdateUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/stats", handler.GetStatsHandler)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"user-service"}`)
	})

	// Start server
	port := ":8081"
	log.Printf("User Service starting on port %s", port)
	log.Println("Available endpoints:")
	log.Println("  POST /register - Register a new user")
	log.Println("  POST /login - User login")
	log.Println("  GET /users/{id} - Get user by ID")
	log.Println("  GET /health - Health check")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
