# Frontend Integration Guide

## CORS Issue Fixed ✅

The CORS error has been resolved by removing duplicate CORS headers from backend services. Only the API Gateway now handles CORS.

## API Endpoints (All via Gateway: http://localhost:8080)

### User Management
```javascript
// Register User
POST /api/users/register
{
  "username": "john_doe",
  "email": "john@example.com", 
  "password": "password123"
}

// Login User
POST /api/users/login
{
  "email": "john@example.com",
  "password": "password123"
}

// Get User by ID
GET /api/users/{id}
```

### Appointment Management
```javascript
// Create Appointment
POST /api/appointments
{
  "user_id": 1,
  "title": "Doctor Appointment",
  "description": "Annual checkup",
  "date": "2024-02-15",
  "time": "14:30",
  "duration": 30
}

// Get Appointments by User
GET /api/appointments?user_id=1

// Get Specific Appointment
GET /api/appointments/{id}

// Update Appointment
PUT /api/appointments/{id}
{
  "user_id": 1,
  "title": "Updated Appointment",
  "description": "Updated description",
  "date": "2024-02-15",
  "time": "15:00",
  "duration": 45
}

// Delete Appointment
DELETE /api/appointments/{id}
```

## Frontend Integration Example

### JavaScript/Fetch API
```javascript
const API_BASE = 'http://localhost:8080';

// Register User
async function registerUser(userData) {
  const response = await fetch(`${API_BASE}/api/users/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData)
  });
  return response.json();
}

// Create Appointment
async function createAppointment(appointmentData) {
  const response = await fetch(`${API_BASE}/api/appointments`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(appointmentData)
  });
  return response.json();
}

// Get User Appointments
async function getUserAppointments(userId) {
  const response = await fetch(`${API_BASE}/api/appointments?user_id=${userId}`);
  return response.json();
}

// Usage Example
registerUser({
  username: 'test_user',
  email: 'test@example.com',
  password: 'password123'
}).then(data => {
  console.log('User registered:', data);
  
  // Create appointment for the user
  return createAppointment({
    user_id: data.user_id,
    title: 'Test Appointment',
    description: 'Testing frontend integration',
    date: '2024-02-20',
    time: '10:00',
    duration: 30
  });
}).then(appointment => {
  console.log('Appointment created:', appointment);
});
```

### React Example
```jsx
import React, { useState, useEffect } from 'react';

const AppointmentApp = () => {
  const [user, setUser] = useState(null);
  const [appointments, setAppointments] = useState([]);

  const registerUser = async (userData) => {
    try {
      const response = await fetch('http://localhost:8080/api/users/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(userData)
      });
      const data = await response.json();
      setUser(data);
      return data;
    } catch (error) {
      console.error('Registration failed:', error);
    }
  };

  const createAppointment = async (appointmentData) => {
    try {
      const response = await fetch('http://localhost:8080/api/appointments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(appointmentData)
      });
      return response.json();
    } catch (error) {
      console.error('Appointment creation failed:', error);
    }
  };

  const loadAppointments = async (userId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/appointments?user_id=${userId}`);
      const data = await response.json();
      setAppointments(data);
    } catch (error) {
      console.error('Failed to load appointments:', error);
    }
  };

  return (
    <div>
      <h1>Appointment Booking System</h1>
      {/* Your UI components here */}
    </div>
  );
};
```

## Response Formats

### Success Response
```json
{
  "user_id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "token": "token_1"
}
```

### Appointment Response
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Doctor Appointment",
  "description": "Annual checkup",
  "date": "2024-02-15",
  "time": "14:30",
  "duration": 30,
  "status": "scheduled",
  "created_at": "2026-02-04T02:08:14Z",
  "updated_at": "2026-02-04T02:08:14Z"
}
```

### Error Response
```json
{
  "error": "Bad Request",
  "message": "Username, email, and password are required"
}
```

## CORS Headers
All responses include:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Testing the Integration

1. Start all services: Run `start-services.bat`
2. Test APIs: Run `test-apis.bat`
3. Your frontend (running on any port) can now access the APIs without CORS errors

## Common Issues & Solutions

### CORS Errors
- ✅ Fixed: Only API Gateway sets CORS headers
- Ensure your frontend uses `http://localhost:8080` as the API base URL

### Connection Refused
- Make sure all three services are running
- Check ports: 8080 (Gateway), 8081 (User), 8082 (Appointment)

### 404 Errors
- Verify URL paths include `/api/` prefix
- Check HTTP methods match the endpoint requirements
