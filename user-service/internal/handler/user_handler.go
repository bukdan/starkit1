package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"user-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterProtectedRoutes registers profile and admin endpoints
func RegisterProtectedRoutes(rg *gin.RouterGroup, db *pgxpool.Pool) {
	repo := repository.NewUserRepository(db)
	v := rg.Group("/users")

	v.GET("/me", AuthMiddleware(), func(c *gin.Context) {
		uid := c.GetString("sub")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		u, err := repo.GetByID(ctx, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	v.PUT("/me", AuthMiddleware(), func(c *gin.Context) {
		uid := c.GetString("sub")
		var in struct {
			Username  string `json:"username"`
			Phone     string `json:"phone"`
			AvatarURL string `json:"avatar_url"`
		}
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		if err := repo.UpdateProfile(ctx, uid, in.Username, in.Phone, in.AvatarURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	})

	// Admin routes
	admin := rg.Group("/admin", AuthMiddleware(), AdminOnly())
	admin.GET("/users", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		users, err := repo.ListUsers(ctx, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	admin.GET("/users/:id", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		u, err := repo.GetByID(ctx, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	admin.PUT("/users/:id", func(c *gin.Context) {
		var in struct {
			Username   string `json:"username"`
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			Role       string `json:"role"`
			IsVerified bool   `json:"is_verified"`
		}
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		_ = repo.UpdateUserAdmin(ctx, c.Param("id"), in.Username, in.Email, in.Phone, in.Role, in.IsVerified)
		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	})

	admin.DELETE("/users/:id", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		_ = repo.DeleteUser(ctx, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})
}
