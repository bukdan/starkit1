r := gin.Default()

// Auth routes
authHandler := handler.NewAuthHandler(authService)
r.POST("/auth/register", authHandler.Register)
r.POST("/auth/login", authHandler.Login)
r.POST("/auth/login-google", authHandler.LoginWithGoogle)
r.POST("/auth/verify-otp", authHandler.VerifyOTP)

// User routes (protected by JWT)
userHandler := handler.NewUserHandler(userService)
auth := r.Group("/user", middleware.JWTAuthMiddleware())
{
    auth.GET("/profile", userHandler.GetProfile)
    auth.PUT("/profile", userHandler.UpdateProfile)
}
==============================

✅ Alur Kerja Middleware JWT
Client kirim request dengan Authorization Header:


Authorization: Bearer <JWT_TOKEN>
Middleware cek validitas token.

Jika valid → inject user_id ke dalam context.

Handler bisa ambil c.GetString("user_id").

Jika token invalid/expired → return 401 Unauthorized.
✅ Testing Endpoint

Register User

POST http://localhost:8080/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123",
  "phone": "628123456789"
}


→ Kirim OTP ke email/WhatsApp.

Verify OTP

POST http://localhost:8080/auth/verify-otp
{
  "email": "john@example.com",
  "otp": "123456"
}


Login

POST http://localhost:8080/auth/login
{
  "email": "john@example.com",
  "password": "secret123"
}


→ Response JWT token.

Profile (Protected)

GET http://localhost:8080/user/profile
Header: Authorization: Bearer <JWT_TOKEN>


🔥 Dengan ini user-service kita sudah siap dijalankan end-to-end.
=================================================
✅ Testing Endpoint

Register User

POST http://localhost:8080/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123",
  "phone": "628123456789"
}


→ Kirim OTP ke email/WhatsApp.

Verify OTP

POST http://localhost:8080/auth/verify-otp
{
  "email": "john@example.com",
  "otp": "123456"
}


Login

POST http://localhost:8080/auth/login
{
  "email": "john@example.com",
  "password": "secret123"
}


→ Response JWT token.

Profile (Protected)

GET http://localhost:8080/user/profile
Header: Authorization: Bearer <JWT_TOKEN>