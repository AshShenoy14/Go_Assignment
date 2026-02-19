package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Initialize database
	db, err := NewDatabase("appointments.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handler
	handler := NewHandler(db)

	// Setup routes
	http.HandleFunc("/appointments", func(w http.ResponseWriter, r *http.Request) {
		// Route based on HTTP method
		switch r.Method {
		case http.MethodPost:
			handler.CreateAppointmentHandler(w, r)
		case http.MethodGet:
			handler.GetAppointmentsHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/appointments/slots", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSlotsHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/appointments/", func(w http.ResponseWriter, r *http.Request) {
		// Handle subpaths
		path := r.URL.Path
		if strings.HasSuffix(path, "/cancel") {
			if r.Method == http.MethodPost {
				handler.CancelAppointmentHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Standard ID based routes
		switch r.Method {
		case http.MethodGet:
			handler.GetAppointmentHandler(w, r)
		case http.MethodPut:
			handler.UpdateAppointmentHandler(w, r)
		case http.MethodDelete:
			handler.DeleteAppointmentHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetNotificationsHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetStatsHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"appointment-service"}`)
	})

	// Start server
	port := ":8082"
	log.Printf("Appointment Service starting on port %s", port)
	log.Println("Available endpoints:")
	log.Println("  POST /appointments - Create a new appointment")
	log.Println("  GET /appointments?user_id={id} - Get all appointments for a user")
	log.Println("  GET /appointments/{id} - Get appointment by ID")
	log.Println("  PUT /appointments/{id} - Update appointment")
	log.Println("  DELETE /appointments/{id} - Delete appointment")
	log.Println("  GET /health - Health check")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
