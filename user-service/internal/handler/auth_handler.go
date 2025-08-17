package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DTOs
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	SendVia  string `json:"send_via"` // "email" or "wa"
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyOTPInput struct {
	UserID  string `json:"user_id" binding:"required"`
	Channel string `json:"channel" binding:"required"`
	Code    string `json:"code" binding:"required"`
}

type ResendOTPInput struct {
	Email   string `json:"email"`
	UserID  string `json:"user_id"`
	Channel string `json:"channel" binding:"required"`
}

// RegisterRoutes registers auth endpoints
func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool, otpExpiryMinutes int) {
	repo := repository.NewUserRepository(db)
	otpSvc := service.NewOTPService(repo)
	userSvc := service.NewUserService(repo, otpSvc)

	auth := rg.Group("/auth")
	auth.POST("/register", func(c *gin.Context) {
		var in RegisterInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		sendVia := strings.ToLower(in.SendVia)
		if sendVia != "wa" {
			sendVia = "email"
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 7*time.Second)
		defer cancel()
		uid, err := userSvc.Register(ctx, in.Username, strings.ToLower(in.Email), in.Phone, in.Password, sendVia, otpExpiryMinutes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"user_id": uid, "message": "user created, OTP sent"})
	})

	auth.POST("/verify-otp", func(c *gin.Context) {
		var in VerifyOTPInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		if err := userSvc.VerifyOTP(ctx, in.UserID, in.Channel, in.Code); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid/expired code"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "verified"})
	})

	auth.POST("/resend-otp", func(c *gin.Context) {
		var in ResendOTPInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		var u repository.UserRepository // not used, keep compile-safe
		_ = u

		// find user first
		repo := repository.NewUserRepository(db)
		var user model.User
		var err error
		if in.UserID != "" {
			user, err = repo.GetByID(ctx, in.UserID)
		} else {
			user, _, err = repo.GetByEmail(ctx, strings.ToLower(in.Email))
		}
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		contact := user.Email
		if in.Channel == "wa" && user.Phone != nil {
			contact = *user.Phone
		}
		_, _ = otpSvc.GenerateCreateAndSend(ctx, user.ID, in.Channel, contact, otpExpiryMinutes)
		c.JSON(http.StatusOK, gin.H{"message": "otp resent"})
	})

	auth.POST("/login", func(c *gin.Context) {
		var in LoginInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		repo := repository.NewUserRepository(db)
		user, pass, err := repo.GetByEmail(ctx, strings.ToLower(in.Email))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if !utils.ComparePassword(pass, in.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		// create token (24h)
		token, err := auth.IssueToken(user.ID, user.Username, user.Email, user.Role, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
	})

	// Google login endpoint (verify token & upsert user)
	auth.POST("/google", func(c *gin.Context) {
		var in struct {
			IDToken string `json:"id_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		prof, err := auth.VerifyGoogleIDToken(ctx, in.IDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid google id token"})
			return
		}
		// upsert by email
		repo := repository.NewUserRepository(db)
		u, _, err := repo.GetByEmail(ctx, strings.ToLower(prof.Email))
		if err != nil {
			// create user with empty password
			uid, cerr := repo.CreateUser(ctx, prof.Name, strings.ToLower(prof.Email), "", "")
			if cerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
			// mark verified
			_ = repo.SetVerified(ctx, uid)
			u, _ = repo.GetByID(ctx, uid)
		}
		// issue JWT
		token, _ := auth.IssueToken(u.ID, u.Username, u.Email, u.Role, 24*time.Hour)
		c.JSON(http.StatusOK, gin.H{"token": token, "user": u})
	})
}
