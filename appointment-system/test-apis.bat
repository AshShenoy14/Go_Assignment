@echo off
echo Testing Appointment Booking APIs...
echo.

echo [1] Testing API Gateway Health...
curl -UseBasicParsing http://localhost:8080/health
echo.

echo [2] Testing User Registration...
curl -UseBasicParsing -Method POST -Uri http://localhost:8080/api/users/register -ContentType "application/json" -Body '{"username":"demo_user","email":"demo@example.com","password":"password123"}'
echo.

echo [3] Testing Appointment Creation...
curl -UseBasicParsing -Method POST -Uri http://localhost:8080/api/appointments -ContentType "application/json" -Body '{"user_id":1,"title":"Demo Appointment","description":"Testing appointment system","date":"2024-02-25","time":"14:00","duration":45}'
echo.

echo [4] Testing Get Appointment...
curl -UseBasicParsing -Uri "http://localhost:8080/api/appointments/1"
echo.

echo [5] Testing Update Appointment...
curl -UseBasicParsing -Method PUT -Uri "http://localhost:8080/api/appointments/1" -ContentType "application/json" -Body '{"user_id":1,"title":"Updated Demo Appointment","description":"Updated test appointment","date":"2024-02-25","time":"15:00","duration":60}'
echo.

echo [6] Testing Delete Appointment...
curl -UseBasicParsing -Method DELETE -Uri "http://localhost:8080/api/appointments/1"
echo.

echo.
echo All API tests completed!
pause
