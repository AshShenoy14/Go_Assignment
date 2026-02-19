@echo off
echo Starting Appointment Booking Microservices...
echo.

echo [1/3] Starting User Service (Port 8081)...
start "User Service" cmd /k "cd /d %~dp0user-service && go run main.go database.go handlers.go models.go"
timeout /t 2 >nul

echo [2/3] Starting Appointment Service (Port 8082)...
start "Appointment Service" cmd /k "cd /d %~dp0appointment-service && go run main.go database.go handlers.go models.go"
timeout /t 2 >nul

echo [3/3] Starting API Gateway (Port 8080)...
start "API Gateway" cmd /k "cd /d %~dp0api-gateway && go run main.go"
timeout /t 2 >nul

echo.
echo All services started successfully!
echo.
echo API Endpoints:
echo   API Gateway: http://localhost:8080
echo   User Service: http://localhost:8081
echo   Appointment Service: http://localhost:8082
echo.
echo Health Checks:
echo   curl http://localhost:8080/health
echo   curl http://localhost:8081/health
echo   curl http://localhost:8082/health
echo.
pause
