package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_super_secret_key_change_in_production")

const (
	userServiceURL        = "http://localhost:8081"
	appointmentServiceURL = "http://localhost:8082"
)

// ProxyRequest represents a request to be forwarded to a service
type ProxyRequest struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string]string
}

// API Gateway acts as a reverse proxy to route requests to appropriate microservices
func main() {
	// Setup routes for the API Gateway
	http.HandleFunc("/", enableCORS(authMiddleware(routeHandler)))

	// Health check endpoint for the gateway
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"api-gateway"}`)
	})

	// Start the API Gateway server
	port := ":8080"
	log.Printf("API Gateway starting on port %s", port)
	log.Println("Available endpoints:")
	log.Println("  User Service:")
	log.Println("    POST /api/users/register - Register a new user")
	log.Println("    POST /api/users/login - User login")
	log.Println("    GET /api/users/{id} - Get user by ID")
	log.Println("  Appointment Service:")
	log.Println("    POST /api/appointments - Create a new appointment")
	log.Println("    GET /api/appointments?user_id={id} - Get all appointments for a user")
	log.Println("    GET /api/appointments/{id} - Get appointment by ID")
	log.Println("    PUT /api/appointments/{id} - Update appointment")
	log.Println("    DELETE /api/appointments/{id} - Delete appointment")
	log.Println("  Gateway:")
	log.Println("    GET /health - Gateway health check")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

// enableCORS sets CORS headers for all responses
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// authMiddleware validates the JWT token and sets user context
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Public endpoints
		if r.URL.Path == "/health" ||
			r.URL.Path == "/" ||
			strings.HasPrefix(r.URL.Path, "/api/users/login") ||
			strings.HasPrefix(r.URL.Path, "/api/users/register") {
			next(w, r)
			return
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and set headers
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				r.Header.Set("X-User-Id", fmt.Sprintf("%d", int(userID)))
			}
			if role, ok := claims["role"].(string); ok {
				r.Header.Set("X-User-Role", role)
			}
		}

		next(w, r)
	}
}

// routeHandler routes incoming requests to the appropriate microservice
func routeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	log.Printf("Incoming request: %s %s", method, path)

	// Route to User Service
	if strings.HasPrefix(path, "/api/users") {
		handleUserServiceRequest(w, r)
		return
	}

	// Route to Appointment Service
	if strings.HasPrefix(path, "/api/appointments") {
		handleAppointmentServiceRequest(w, r)
		return
	}

	// Route to Notifications (in Appointment Service)
	if strings.HasPrefix(path, "/api/notifications") {
		handleNotificationServiceRequest(w, r)
		return
	}

	// Handle health check
	if path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"api-gateway"}`)
		return
	}

	// Handle root path
	if path == "/" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"Appointment Booking API Gateway","version":"1.0.0","services":["user-service","appointment-service"]}`)
		return
	}

	// 404 for unknown paths
	http.Error(w, "Endpoint not found", http.StatusNotFound)
}

// handleUserServiceRequest forwards requests to the User Service
func handleUserServiceRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Remove /api/users prefix and map to user service endpoints
	var targetPath string
	if strings.HasSuffix(path, "/register") {
		targetPath = "/register"
	} else if strings.HasSuffix(path, "/login") {
		targetPath = "/login"
	} else if strings.HasSuffix(path, "/stats") {
		targetPath = "/stats"
	} else {
		// Handle /api/users/{id} or /api/users (list)
		parts := strings.Split(path, "/")
		// parts[0] = "", parts[1] = "api", parts[2] = "users"
		if len(parts) == 3 || (len(parts) == 4 && parts[3] == "") {
			// List users: /api/users or /api/users/
			targetPath = "/users"
		} else if len(parts) >= 4 && parts[3] != "" {
			// Get User by ID: /api/users/{id}
			targetPath = "/users/" + parts[3]
		} else {
			http.Error(w, "Invalid user endpoint", http.StatusBadRequest)
			return
		}
	}

	// Forward request to user service
	if r.URL.RawQuery != "" {
		targetPath += "?" + r.URL.RawQuery
	}
	proxyRequest(w, r, userServiceURL+targetPath)
}

// handleAppointmentServiceRequest forwards requests to the Appointment Service
func handleAppointmentServiceRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Remove /api prefix and keep the rest
	targetPath := strings.TrimPrefix(path, "/api")

	// Forward request to appointment service
	if r.URL.RawQuery != "" {
		targetPath += "?" + r.URL.RawQuery
	}
	proxyRequest(w, r, appointmentServiceURL+targetPath)
}

// handleNotificationServiceRequest forwards requests to the Appointment Service (where notifications live)
func handleNotificationServiceRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// Remove /api prefix
	targetPath := strings.TrimPrefix(path, "/api")

	if r.URL.RawQuery != "" {
		targetPath += "?" + r.URL.RawQuery
	}
	proxyRequest(w, r, appointmentServiceURL+targetPath)
}

// proxyRequest forwards the HTTP request to the target service
func proxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	// Read the request body
	var body io.Reader
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		body = bytes.NewReader(bodyBytes)
		r.Body = io.NopCloser(body) // Restore body for potential reuse
	}

	// Create new request to target service
	req, err := http.NewRequest(r.Method, targetURL, body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make the request to the target service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error forwarding request to %s: %v", targetURL, err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
		return
	}

	log.Printf("Forwarded %s %s -> %s (status: %d)", r.Method, r.URL.Path, targetURL, resp.StatusCode)
}
