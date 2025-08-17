package main

import (
	"user-service/config"
	"user-service/handler"
	"user-service/repository"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	repo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(repo)
	authHandler := handler.NewAuthHandler(userService)

	r := gin.Default()

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	r.Run(":8081")
}
