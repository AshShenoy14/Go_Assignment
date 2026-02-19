# Appointment Booking and Management System

A simple microservices-based appointment booking system built with Go. This project demonstrates microservices architecture for a college assignment.

## Architecture

The system consists of 3 microservices:

1. **API Gateway** (Port 8080) - Entry point for all client requests
2. **User Service** (Port 8081) - Handles user registration and authentication
3. **Appointment Service** (Port 8082) - Manages appointment operations

## Project Structure

```
appointment-system/
│
├── api-gateway/
│ └── main.go
│
├── user-service/
│ ├── main.go
│ ├── database.go
│ ├── handlers.go
│ └── models.go
│
├── appointment-service/
│ ├── main.go
│ ├── database.go
│ ├── handlers.go
│ └── models.go
│
└── README.md
```

## Prerequisites

- Go 1.19 or higher
- SQLite3 driver for Go

## Setup Instructions

### 1. Install Dependencies

Install the SQLite3 driver:

```bash
go get github.com/mattn/go-sqlite3
go get golang.org/x/crypto/bcrypt
```

### 2. Start the Services

Open 3 separate terminal windows and start each service:

**Terminal 1 - User Service:**
```bash
cd user-service
go run main.go
```

**Terminal 2 - Appointment Service:**
```bash
cd appointment-service
go run main.go
```

**Terminal 3 - API Gateway:**
```bash
cd api-gateway
go run main.go
```

The services will start on the following ports:
- API Gateway: http://localhost:8080
- User Service: http://localhost:8081
- Appointment Service: http://localhost:8082

### 3. Verify Services are Running

Check health endpoints:

```bash
# API Gateway
curl http://localhost:8080/health

# User Service
curl http://localhost:8081/health

# Appointment Service
curl http://localhost:8082/health
```

## API Endpoints

All requests should be made to the **API Gateway** at `http://localhost:8080`

### User Service Endpoints

#### Register User
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

#### Login User
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

#### Get User by ID
```bash
curl -X GET http://localhost:8080/api/users/1
```

### Appointment Service Endpoints

#### Create Appointment
```bash
curl -X POST http://localhost:8080/api/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "Doctor Appointment",
    "description": "Annual checkup",
    "date": "2024-02-15",
    "time": "14:30",
    "duration": 30
  }'
```

#### Get All Appointments for User
```bash
curl -X GET "http://localhost:8080/api/appointments?user_id=1"
```

#### Get Specific Appointment
```bash
curl -X GET http://localhost:8080/api/appointments/1
```

#### Update Appointment
```bash
curl -X PUT http://localhost:8080/api/appointments/1 \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "Doctor Appointment (Updated)",
    "description": "Annual checkup - updated time",
    "date": "2024-02-15",
    "time": "15:00",
    "duration": 45
  }'
```

#### Delete Appointment
```bash
curl -X DELETE http://localhost:8080/api/appointments/1
```

## How Microservices Communicate

1. **Client → API Gateway**: All client requests go to the API Gateway (port 8080)
2. **API Gateway → Services**: The Gateway forwards requests to the appropriate microservice:
   - `/api/users/*` → User Service (port 8081)
   - `/api/appointments/*` → Appointment Service (port 8082)
3. **Service → Database**: Each service has its own SQLite database:
   - User Service: `users.db`
   - Appointment Service: `appointments.db`

## Features

### User Service
- ✅ User registration with password hashing
- ✅ User login with password verification
- ✅ Get user information by ID
- ✅ Email and username uniqueness validation

### Appointment Service
- ✅ Create appointments with conflict prevention
- ✅ View appointments by user ID
- ✅ Update existing appointments
- ✅ Delete appointments
- ✅ Prevent double booking (same user, date, and time)
- ✅ Date and time validation

### API Gateway
- ✅ Request routing to appropriate services
- ✅ CORS support for web applications
- ✅ Request/response forwarding
- ✅ Health check endpoint

## Database Schema

### Users Table (users.db)
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Appointments Table (appointments.db)
```sql
CREATE TABLE appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    date TEXT NOT NULL,
    time TEXT NOT NULL,
    duration INTEGER NOT NULL DEFAULT 30,
    status TEXT NOT NULL DEFAULT 'scheduled',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE(user_id, date, time)
);
```

## Error Handling

The system returns appropriate HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid credentials)
- `404` - Not Found
- `409` - Conflict (duplicate booking)
- `500` - Internal Server Error
- `503` - Service Unavailable

## Security Features

- Password hashing using bcrypt
- Input validation for all endpoints
- SQL injection prevention using parameterized queries
- CORS headers for cross-origin requests

## Testing the System

1. **Register a user** using the registration endpoint
2. **Login** to get the user ID
3. **Create appointments** for the user
4. **Try to create a duplicate appointment** (should fail with conflict error)
5. **Update and delete appointments** to test full CRUD operations

## Notes

- This is a simple demonstration project for educational purposes
- In production, you would add:
  - JWT tokens for authentication
  - Service discovery
  - Message brokers for async communication
  - Container orchestration (Docker/Kubernetes)
  - More robust error handling and logging
  - Database migrations
  - API documentation (Swagger/OpenAPI)

## Troubleshooting

### Port Already in Use
If you get a "port already in use" error, change the port in the respective `main.go` file.

### Database Connection Issues
Ensure SQLite3 driver is installed:
```bash
go get github.com/mattn/go-sqlite3
```

### Service Unavailable
Make sure all services are running before making API calls through the gateway.
