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

	// Create users table
	if err := database.createUsersTable(); err != nil {
		return nil, fmt.Errorf("failed to create users table: %v", err)
	}

	// Migrations
	// Add role column
	db.Exec("ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'patient'")
	// Add specialization column
	db.Exec("ALTER TABLE users ADD COLUMN specialization TEXT DEFAULT ''")
	// Add phone column
	db.Exec("ALTER TABLE users ADD COLUMN phone TEXT DEFAULT ''")

	return database, nil
}

// createUsersTable creates the users table if it doesn't exist
func (d *Database) createUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'patient',
		specialization TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	log.Println("Users table created or already exists")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// CreateUser inserts a new user into the database
func (d *Database) CreateUser(username, email, hashedPassword, role, specialization, phone string) (int, error) {
	query := `INSERT INTO users (username, email, password, role, specialization, phone) VALUES (?, ?, ?, ?, ?, ?)`

	// Default to patient if empty
	if role == "" {
		role = "patient"
	}

	result, err := d.db.Exec(query, username, email, hashedPassword, role, specialization, phone)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %v", err)
	}

	return int(userID), nil
}

// GetUserByEmail retrieves a user by email
func (d *Database) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password, role, specialization, phone, created_at FROM users WHERE email = ?`

	row := d.db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.Specialization, &user.Phone, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(userID int) (*User, error) {
	query := `SELECT id, username, email, role, specialization, phone, created_at FROM users WHERE id = ?`

	row := d.db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Specialization, &user.Phone, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// GetUsersByRole retrieves all users with a specific role
func (d *Database) GetUsersByRole(role string) ([]*User, error) {
	query := `SELECT id, username, email, role, specialization, phone FROM users WHERE role = ?`

	rows, err := d.db.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Specialization, &user.Phone); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

// UpdateUser updates an existing user
func (d *Database) UpdateUser(userID int, username, email, password, specialization, phone string) error {
	var query string
	var args []interface{}

	if password != "" {
		query = `UPDATE users SET username = ?, email = ?, password = ?, specialization = ?, phone = ? WHERE id = ?`
		args = []interface{}{username, email, password, specialization, phone, userID}
	} else {
		query = `UPDATE users SET username = ?, email = ?, specialization = ?, phone = ? WHERE id = ?`
		args = []interface{}{username, email, specialization, phone, userID}
	}

	_, err := d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

// GetStats returns the total number of users and doctors
func (d *Database) GetStats() (map[string]int, error) {
	stats := make(map[string]int)

	var totalUsers int
	err := d.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		return nil, err
	}
	stats["total_users"] = totalUsers

	var totalDoctors int
	err = d.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'doctor'").Scan(&totalDoctors)
	if err != nil {
		return nil, err
	}
	stats["total_doctors"] = totalDoctors

	return stats, nil
}
