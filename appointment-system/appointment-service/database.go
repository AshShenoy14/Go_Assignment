package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection and initializes tables
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	database := &Database{db: db}

	// Create appointments table
	if err := database.createAppointmentsTable(); err != nil {
		return nil, fmt.Errorf("failed to create appointments table: %v", err)
	}

	// Create availability table
	if err := database.createAvailabilityTable(); err != nil {
		return nil, fmt.Errorf("failed to create availability table: %v", err)
	}

	// Create notifications table
	if err := database.createNotificationsTable(); err != nil {
		return nil, fmt.Errorf("failed to create notifications table: %v", err)
	}

	// Migration: Add doctor_id column if it doesn't exist
	db.Exec("ALTER TABLE appointments ADD COLUMN doctor_id INTEGER NOT NULL DEFAULT 0")
	// Add cancelled_reason column
	db.Exec("ALTER TABLE appointments ADD COLUMN cancelled_reason TEXT DEFAULT ''")

	return database, nil
}

// createAppointmentsTable creates the appointments table if it doesn't exist
func (d *Database) createAppointmentsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS appointments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		doctor_id INTEGER NOT NULL DEFAULT 0,
		title TEXT NOT NULL,
		description TEXT,
		date TEXT NOT NULL,
		time TEXT NOT NULL,
		duration INTEGER NOT NULL DEFAULT 30,
		status TEXT NOT NULL DEFAULT 'scheduled',
		cancelled_reason TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
		UNIQUE(user_id, date, time)
	)`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create appointments table: %v", err)
	}

	log.Println("Appointments table created or already exists")
	return nil
}

// createNotificationsTable creates the notifications table if it doesn't exist
func (d *Database) createNotificationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		message TEXT NOT NULL,
		type TEXT NOT NULL, -- 'booking', 'cancellation', 'reschedule'
		is_read BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create notifications table: %v", err)
	}

	log.Println("Notifications table created or already exists")
	return nil
}

// CreateNotification creates a new notification
func (d *Database) CreateNotification(userID int, message, notifType string) error {
	query := `INSERT INTO notifications (user_id, message, type) VALUES (?, ?, ?)`
	_, err := d.db.Exec(query, userID, message, notifType)
	return err
}

// GetNotificationsByUserID retrieves notifications for a user
func (d *Database) GetNotificationsByUserID(userID int) ([]*Notification, error) {
	query := `SELECT id, user_id, message, type, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

// MarkNotificationAsRead marks a notification as read
func (d *Database) MarkNotificationAsRead(id int) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = ?`
	_, err := d.db.Exec(query, id)
	return err
}

// createAvailabilityTable creates the availability table if it doesn't exist
func (d *Database) createAvailabilityTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS availability (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		doctor_id INTEGER NOT NULL,
		day_of_week INTEGER NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		UNIQUE(doctor_id, day_of_week)
	)`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create availability table: %v", err)
	}

	log.Println("Availability table created or already exists")
	return nil
}

// SetAvailability sets or updates a doctor's availability for a day
func (d *Database) SetAvailability(doctorID, dayOfWeek int, startTime, endTime string) error {
	query := `INSERT INTO availability (doctor_id, day_of_week, start_time, end_time)
			  VALUES (?, ?, ?, ?)
			  ON CONFLICT(doctor_id, day_of_week) 
			  DO UPDATE SET start_time=excluded.start_time, end_time=excluded.end_time`

	_, err := d.db.Exec(query, doctorID, dayOfWeek, startTime, endTime)
	if err != nil {
		return fmt.Errorf("failed to set availability: %v", err)
	}
	return nil
}

// GetAvailability retrieves a doctor's availability for a specific day
func (d *Database) GetAvailability(doctorID, dayOfWeek int) (*Availability, error) {
	query := `SELECT id, doctor_id, day_of_week, start_time, end_time 
			  FROM availability WHERE doctor_id = ? AND day_of_week = ?`

	row := d.db.QueryRow(query, doctorID, dayOfWeek)

	var av Availability
	err := row.Scan(&av.ID, &av.DoctorID, &av.DayOfWeek, &av.StartTime, &av.EndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no availability set (not an error)
		}
		return nil, fmt.Errorf("failed to get availability: %v", err)
	}

	return &av, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// CreateAppointment inserts a new appointment into the database
func (d *Database) CreateAppointment(appointment *AppointmentRequest) (int, error) {
	query := `INSERT INTO appointments (user_id, doctor_id, title, description, date, time, duration, status) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, 'scheduled')`

	result, err := d.db.Exec(query, appointment.UserID, appointment.DoctorID, appointment.Title,
		appointment.Description, appointment.Date, appointment.Time, appointment.Duration)
	if err != nil {
		return 0, fmt.Errorf("failed to create appointment: %v", err)
	}

	appointmentID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get appointment ID: %v", err)
	}

	return int(appointmentID), nil
}

// GetAppointmentByID retrieves an appointment by ID
func (d *Database) GetAppointmentByID(appointmentID int) (*Appointment, error) {
	query := `SELECT id, user_id, doctor_id, title, description, date, time, duration, status, 
			  created_at, updated_at FROM appointments WHERE id = ?`

	row := d.db.QueryRow(query, appointmentID)

	var appointment Appointment
	err := row.Scan(&appointment.ID, &appointment.UserID, &appointment.DoctorID, &appointment.Title,
		&appointment.Description, &appointment.Date, &appointment.Time,
		&appointment.Duration, &appointment.Status, &appointment.CreatedAt, &appointment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("failed to get appointment: %v", err)
	}

	return &appointment, nil
}

// GetAppointmentsByUserID retrieves all appointments for a specific user
func (d *Database) GetAppointmentsByUserID(userID int) ([]*Appointment, error) {
	query := `SELECT id, user_id, doctor_id, title, description, date, time, duration, status, 
			  created_at, updated_at FROM appointments WHERE user_id = ? ORDER BY date, time`

	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %v", err)
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		err := rows.Scan(&appointment.ID, &appointment.UserID, &appointment.DoctorID, &appointment.Title,
			&appointment.Description, &appointment.Date, &appointment.Time,
			&appointment.Duration, &appointment.Status, &appointment.CreatedAt, &appointment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %v", err)
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

// GetAppointmentsByDoctorID retrieves all appointments for a specific doctor on a specific date
// GetAppointmentsByDoctorID retrieves all appointments for a specific doctor on a specific date
func (d *Database) GetAppointmentsByDoctorID(doctorID int, date string) ([]*Appointment, error) {
	query := `SELECT id, user_id, doctor_id, title, description, date, time, duration, status, 
			  created_at, updated_at FROM appointments WHERE doctor_id = ? AND date = ? ORDER BY time`

	rows, err := d.db.Query(query, doctorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %v", err)
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		err := rows.Scan(&appointment.ID, &appointment.UserID, &appointment.DoctorID, &appointment.Title,
			&appointment.Description, &appointment.Date, &appointment.Time,
			&appointment.Duration, &appointment.Status, &appointment.CreatedAt, &appointment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %v", err)
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

// UpdateAppointment updates an existing appointment
func (d *Database) UpdateAppointment(appointmentID int, appointment *AppointmentRequest) error {
	query := `UPDATE appointments SET title = ?, description = ?, date = ?, 
			  time = ?, duration = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := d.db.Exec(query, appointment.Title, appointment.Description,
		appointment.Date, appointment.Time, appointment.Duration, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to update appointment: %v", err)
	}

	return nil
}

// DeleteAppointment cancels an appointment (Soft Delete)
func (d *Database) DeleteAppointment(appointmentID int, reason string) error {
	query := `UPDATE appointments SET status = 'cancelled', cancelled_reason = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := d.db.Exec(query, reason, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %v", err)
	}

	return nil
}

// CheckAppointmentConflict checks if there's a conflicting appointment
func (d *Database) CheckAppointmentConflict(userID int, date, time string, excludeID int) (bool, error) {
	var query string
	var args []interface{}

	if excludeID > 0 {
		// For updates, exclude the current appointment
		query = `SELECT COUNT(*) FROM appointments 
				 WHERE user_id = ? AND date = ? AND time = ? AND id != ? AND status != 'cancelled'`
		args = []interface{}{userID, date, time, excludeID}
	} else {
		// For new appointments
		query = `SELECT COUNT(*) FROM appointments 
				 WHERE user_id = ? AND date = ? AND time = ? AND status != 'cancelled'`
		args = []interface{}{userID, date, time}
	}

	var count int
	err := d.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check appointment conflict: %v", err)
	}

	return count > 0, nil
}

// GetStats returns appointment statistics
func (d *Database) GetStats() (map[string]int, error) {
	stats := make(map[string]int)

	var totalAppointments int
	d.db.QueryRow("SELECT COUNT(*) FROM appointments").Scan(&totalAppointments)
	stats["total_appointments"] = totalAppointments

	var totalCancelled int
	d.db.QueryRow("SELECT COUNT(*) FROM appointments WHERE status = 'cancelled'").Scan(&totalCancelled)
	stats["cancelled_appointments"] = totalCancelled

	var totalCompleted int
	d.db.QueryRow("SELECT COUNT(*) FROM appointments WHERE status = 'completed'").Scan(&totalCompleted)
	stats["completed_appointments"] = totalCompleted

	return stats, nil
}
