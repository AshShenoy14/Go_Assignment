package main

import "time"

// Appointment represents an appointment in the system
type Appointment struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	DoctorID        int       `json:"doctor_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Date            string    `json:"date"`     // Format: YYYY-MM-DD
	Time            string    `json:"time"`     // Format: HH:MM
	Duration        int       `json:"duration"` // Duration in minutes
	Status          string    `json:"status"`   // "scheduled", "completed", "cancelled"
	CancelledReason string    `json:"cancelled_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Notification represents a system notification
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // 'booking', 'cancellation', 'reschedule'
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// Availability represents a doctor's working hours
type Availability struct {
	ID        int    `json:"id"`
	DoctorID  int    `json:"doctor_id"`
	DayOfWeek int    `json:"day_of_week"` // 0=Sunday, 1=Monday, etc.
	StartTime string `json:"start_time"`  // HH:MM
	EndTime   string `json:"end_time"`    // HH:MM
}

// AppointmentRequest represents appointment creation/update request
type AppointmentRequest struct {
	UserID      int    `json:"user_id"`
	DoctorID    int    `json:"doctor_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`     // Format: YYYY-MM-DD
	Time        string `json:"time"`     // Format: HH:MM
	Duration    int    `json:"duration"` // Duration in minutes
}

// AvailabilityRequest represents request to set availability
type AvailabilityRequest struct {
	DoctorID  int    `json:"doctor_id"`
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// AppointmentResponse represents appointment response with additional info
type AppointmentResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Date            string    `json:"date"`
	Time            string    `json:"time"`
	Duration        int       `json:"duration"`
	Status          string    `json:"status"`
	CancelledReason string    `json:"cancelled_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
