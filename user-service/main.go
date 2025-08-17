package main

import (
	"log"
	"net/http"
	"user-service/config"
	"user-service/handler"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using OS environment variables")
	}

	cfg := config.LoadConfig()

	// Routes
	http.HandleFunc("/auth/register", handler.RegisterUser)
	http.HandleFunc("/auth/login", handler.LoginUser)
	http.HandleFunc("/auth/me", handler.JWTMiddleware(handler.Me))
	http.HandleFunc("/users", handler.JWTMiddleware(handler.GetUsers))

	log.Printf("User-Service running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
