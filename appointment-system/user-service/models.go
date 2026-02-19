package main

import "time"

// User represents a user in the system
// User represents a user in the system
type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"-"` // Don't include password in JSON responses
	Role           string    `json:"role"`
	Specialization string    `json:"specialization"`
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"created_at"`
}

// UserRequest represents user registration/login request
type UserRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`           // Optional, defaults to "patient"
	Specialization string `json:"specialization"` // Optional, for doctors
	Phone          string `json:"phone"`          // Optional
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}
