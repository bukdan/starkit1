package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Routes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "user-service running"})
    })

    // Dummy user handler
    r.GET("/users", GetUsers)

    log.Println("Starting user-service on :8081")
    if err := r.Run(":8081"); err != nil {
        log.Fatal(err)
    }
}
