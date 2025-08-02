@echo off
REM エラーコードテスト用バッチファイル（Windows用）
set BASE_URL=http://localhost:8080

echo === API Error Code Testing ===
echo.

echo 1. Testing validation errors (signup with invalid data):
curl -s -X POST %BASE_URL%/api/auth/signup -H "Content-Type: application/json" -d "{\"email\":\"invalid-email\",\"password\":\"123\",\"display_name\":\"\"}"
echo.

echo 2. Testing validation errors (login with empty data):
curl -s -X POST %BASE_URL%/api/auth/login -H "Content-Type: application/json" -d "{\"email\":\"\",\"password\":\"\"}"
echo.

echo 3. Testing invalid request (malformed JSON):
curl -s -X POST %BASE_URL%/api/auth/signup -H "Content-Type: application/json" -d "{\"email\":\"user@example.com\",\"password\""
echo.

echo 4. Testing successful signup:
curl -s -X POST %BASE_URL%/api/auth/signup -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password123\",\"display_name\":\"Test User\"}"
echo.

echo 5. Testing user exists error (duplicate signup):
curl -s -X POST %BASE_URL%/api/auth/signup -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password123\",\"display_name\":\"Test User\"}"
echo.

echo 6. Testing successful login:
curl -s -X POST %BASE_URL%/api/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
echo.

echo 7. Testing invalid credentials:
curl -s -X POST %BASE_URL%/api/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"wrongpassword\"}"
echo.

echo === Test Complete ===
pause