package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handler contains the database and handles HTTP requests
type Handler struct {
	db *Database
}

// NewHandler creates a new handler with database connection
func NewHandler(db *Database) *Handler {
	return &Handler{db: db}
}

// enableCORS sets CORS headers for all responses
func (h *Handler) enableCORS(next http.HandlerFunc) http.HandlerFunc {
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

// CreateAppointmentHandler handles appointment creation
func (h *Handler) CreateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := h.validateAppointmentRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Override user_id from header (set by API Gateway)
	userIDStr := r.Header.Get("X-User-Id")
	if userIDStr != "" {
		if uid, err := strconv.Atoi(userIDStr); err == nil {
			req.UserID = uid
		}
	} else {
		// If running without gateway (e.g. tests), we might trust the body,
		// but in production, this should be rejected if not internal.
		// For now, we'll allow body if header is missing for easier testing w/o gateway,
		// but ideally gateway is mandatory.
	}

	// Check for appointment conflicts
	hasConflict, err := h.db.CheckAppointmentConflict(req.UserID, req.Date, req.Time, 0)
	if err != nil {
		log.Printf("Error checking appointment conflict: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if hasConflict {
		http.Error(w, "Appointment already exists at this date and time", http.StatusConflict)
		return
	}

	// Create appointment in database
	appointmentID, err := h.db.CreateAppointment(&req)
	if err != nil {
		log.Printf("Error creating appointment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create Notification for the user
	h.db.CreateNotification(req.UserID, fmt.Sprintf("Appointment '%s' confirmed for %s at %s", req.Title, req.Date, req.Time), "booking")

	// Get the created appointment
	appointment, err := h.db.GetAppointmentByID(appointmentID)
	if err != nil {
		log.Printf("Error retrieving created appointment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appointment)
}

// GetAppointmentsHandler handles getting appointments by user ID
func (h *Handler) GetAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from URL query parameter
	userIDStr := r.URL.Query().Get("user_id")
	doctorIDStr := r.URL.Query().Get("doctor_id")

	if userIDStr == "" && doctorIDStr == "" {
		http.Error(w, "User ID or Doctor ID is required", http.StatusBadRequest)
		return
	}

	var appointments []*Appointment
	var err error

	if userIDStr != "" {
		userID, _ := strconv.Atoi(userIDStr)
		// Security Check: Ensure patients can only view their own appointments
		authUserIDStr := r.Header.Get("X-User-Id")
		authRole := r.Header.Get("X-User-Role")

		if authUserIDStr != "" && authRole != "admin" && authRole != "doctor" {
			authUserID, _ := strconv.Atoi(authUserIDStr)
			if authUserID != userID {
				http.Error(w, "Unauthorized to view these appointments", http.StatusForbidden)
				return
			}
		}
		appointments, err = h.db.GetAppointmentsByUserID(userID)
	} else if doctorIDStr != "" {
		doctorID, _ := strconv.Atoi(doctorIDStr)
		date := r.URL.Query().Get("date")
		if date == "" {
			// For doc dashboard, maybe we want all dates? But specialized function GetAppointmentsByDoctorID needs date.
			// Let's modify GetAppointmentsByDoctorID or add a new one in DB for all dates.
			// For now, let's just err if date is missing for doctor view, or implement filtering in memory (bad),
			// or better: The requirement says "View daily appointments", so date is expected.
			// But for "Upcoming and past", we need a range.
			// Let's check DB support. We only have GetAppointmentsByDoctorID(id, date).
			// We might need GetAppointmentsByDoctorIDAll(id).
			// For now, let's require date for doctor view to match "View daily appointments".
			http.Error(w, "Date is required for doctor appointments", http.StatusBadRequest)
			return
		}
		appointments, err = h.db.GetAppointmentsByDoctorID(doctorID, date)
	}

	if err != nil {
		log.Printf("Error getting appointments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// GetAppointmentHandler handles getting a single appointment by ID
func (h *Handler) GetAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get appointment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/appointments/")
	if path == "" {
		http.Error(w, "Appointment ID is required", http.StatusBadRequest)
		return
	}

	appointmentID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	// Get appointment from database
	appointment, err := h.db.GetAppointmentByID(appointmentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Appointment not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting appointment: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointment)
}

// UpdateAppointmentHandler handles appointment updates
func (h *Handler) UpdateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept PUT requests
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get appointment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/appointments/")
	if path == "" {
		http.Error(w, "Appointment ID is required", http.StatusBadRequest)
		return
	}

	appointmentID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	var req AppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := h.validateAppointmentRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if appointment exists
	existingDetails, err := h.db.GetAppointmentByID(appointmentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Appointment not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting appointment: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Check for appointment conflicts (exclude current appointment)
	hasConflict, err := h.db.CheckAppointmentConflict(req.UserID, req.Date, req.Time, appointmentID)
	if err != nil {
		log.Printf("Error checking appointment conflict: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if hasConflict {
		http.Error(w, "Appointment already exists at this date and time", http.StatusConflict)
		return
	}

	// Update appointment in database
	err = h.db.UpdateAppointment(appointmentID, &req)
	if err != nil {
		log.Printf("Error updating appointment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Notify rescheduled
	if existingDetails.Date != req.Date || existingDetails.Time != req.Time {
		h.db.CreateNotification(req.UserID, fmt.Sprintf("Appointment rescheduled to %s at %s", req.Date, req.Time), "reschedule")
	}

	// Get the updated appointment
	updatedAppointment, err := h.db.GetAppointmentByID(appointmentID)
	if err != nil {
		log.Printf("Error retrieving updated appointment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAppointment)
}

// CancelAppointmentHandler handles appointment cancellation (Soft Delete)
func (h *Handler) CancelAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	// Accept POST requests to /appointments/{id}/cancel
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/appointments/")
	path = strings.TrimSuffix(path, "/cancel")
	if path == "" {
		http.Error(w, "Appointment ID is required", http.StatusBadRequest)
		return
	}

	appointmentID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	// Decode body if present, ignore error if empty
	json.NewDecoder(r.Body).Decode(&req)

	// Check if appointment exists
	app, err := h.db.GetAppointmentByID(appointmentID)
	if err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	// Perform soft cancel
	err = h.db.DeleteAppointment(appointmentID, req.Reason)
	if err != nil {
		log.Printf("Error cancelling appointment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create notification
	h.db.CreateNotification(app.UserID, fmt.Sprintf("Appointment '%s' cancelled", app.Title), "cancellation")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Appointment cancelled successfully"})
}

// DeleteAppointmentHandler handles appointment deletion (Soft Cancel via DELETE method)
func (h *Handler) DeleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/appointments/")
	if path == "" {
		http.Error(w, "Appointment ID is required", http.StatusBadRequest)
		return
	}

	appointmentID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	reason := r.URL.Query().Get("reason")
	if reason == "" {
		reason = "Deleted by user"
	}

	// Check if appointment exists
	app, err := h.db.GetAppointmentByID(appointmentID)
	if err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	// Perform soft cancel
	err = h.db.DeleteAppointment(appointmentID, reason)
	if err != nil {
		log.Printf("Error deleting appointment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create notification
	h.db.CreateNotification(app.UserID, fmt.Sprintf("Appointment '%s' deleted", app.Title), "cancellation")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Appointment deleted successfully"})
}

// GetStatsHandler returns appointment statistics
func (h *Handler) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.db.GetStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetNotificationsHandler
func (h *Handler) GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	notifs, err := h.db.GetNotificationsByUserID(userID)
	if err != nil {
		http.Error(w, "Error fetching notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

// validateAppointmentRequest validates the appointment request
func (h *Handler) validateAppointmentRequest(req *AppointmentRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("valid user ID is required")
	}
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Date == "" {
		return fmt.Errorf("date is required")
	}
	if req.Time == "" {
		return fmt.Errorf("time is required")
	}
	if req.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}

	// Validate date format (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	// Validate time format (HH:MM)
	if _, err := time.Parse("15:04", req.Time); err != nil {
		return fmt.Errorf("invalid time format, expected HH:MM")
	}

	return nil
}

// SetAvailabilityHandler handles setting doctor availability
func (h *Handler) SetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate doctor_id (simple check)
	if req.DoctorID <= 0 {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	// Save to DB
	err := h.db.SetAvailability(req.DoctorID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		log.Printf("Error setting availability: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Availability updated successfully"})
}

// GetAvailabilityHandler retrieves doctor availability
func (h *Handler) GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	doctorIDStr := r.URL.Query().Get("doctor_id")
	dayOfWeekStr := r.URL.Query().Get("day_of_week")

	if doctorIDStr == "" || dayOfWeekStr == "" {
		http.Error(w, "doctor_id and day_of_week are required", http.StatusBadRequest)
		return
	}

	doctorID, _ := strconv.Atoi(doctorIDStr)
	dayOfWeek, _ := strconv.Atoi(dayOfWeekStr)

	av, err := h.db.GetAvailability(doctorID, dayOfWeek)
	if err != nil {
		log.Printf("Error getting availability: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if av == nil {
		// Default availability? Or 404? Let's return default 9-5 for now or empty
		// Returning 404 might be better for "not set"
		http.Error(w, "Availability not set", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(av)
}

// GetSlotsHandler generates available time slots
func (h *Handler) GetSlotsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	doctorIDStr := r.URL.Query().Get("doctor_id")
	dateStr := r.URL.Query().Get("date")

	if doctorIDStr == "" || dateStr == "" {
		http.Error(w, "doctor_id and date are required", http.StatusBadRequest)
		return
	}

	doctorID, _ := strconv.Atoi(doctorIDStr)

	// Parse date to get Day of Week
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	dayOfWeek := int(date.Weekday())

	// 1. Get Doctor's Availability
	av, err := h.db.GetAvailability(doctorID, dayOfWeek)
	if err != nil {
		log.Printf("Error getting availability: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Default availability if not set: 9 AM to 5 PM
	startTime := "09:00"
	endTime := "17:00"
	if av != nil {
		startTime = av.StartTime
		endTime = av.EndTime
	}

	// 2. Get Existing Appointments
	appointments, err := h.db.GetAppointmentsByDoctorID(doctorID, dateStr)
	if err != nil {
		log.Printf("Error getting appointments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 3. Generate Slots
	slots := generateSlots(startTime, endTime, 30, appointments)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

// Helper to generate slots
func generateSlots(start, end string, durationMin int, appointments []*Appointment) []string {
	var slots []string

	current, _ := time.Parse("15:04", start)
	endT, _ := time.Parse("15:04", end)

	// Create a map of booked times for O(1) lookup
	booked := make(map[string]bool)
	for _, app := range appointments {
		booked[app.Time] = true
		// Also mark duration? Complex...
		// For now simple Assumption: appointments are also 30 mins aligned
	}

	for current.Before(endT) {
		timeStr := current.Format("15:04")

		if !booked[timeStr] {
			slots = append(slots, timeStr)
		}

		current = current.Add(time.Duration(durationMin) * time.Minute)
	}

	return slots
}
