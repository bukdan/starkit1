package main

import (
	"log"
	"user-service/config"
	"user-service/handler"
	"user-service/middleware"
	"user-service/repository"
	"user-service/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config dari .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Koneksi ke DB PostgreSQL
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Repository
	userRepo := repository.NewUserRepository(db)

	// Services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Gin router
	r := gin.Default()

	// Auth routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/login-google", authHandler.LoginWithGoogle)
	r.POST("/auth/verify-otp", authHandler.VerifyOTP)

	// User routes (protected)
	auth := r.Group("/user", middleware.JWTAuthMiddleware())
	{
		auth.GET("/profile", userHandler.GetProfile)
		auth.PUT("/profile", userHandler.UpdateProfile)
	}

	// Run server
	log.Printf("Server running on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
