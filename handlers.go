
package main

import (
    "context"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    "golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Phone    string `json:"phone"`
    SendVia  string `json:"send_via"`
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

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
    repo := NewUserRepo(db)
    api := r.Group("/api/v1")
    auth := api.Group("/auth")
    auth.POST("/register", func(c *gin.Context){ handleRegister(c, repo) })
    auth.POST("/login", func(c *gin.Context){ handleLogin(c, repo) })
    auth.POST("/verify-otp", func(c *gin.Context){ handleVerifyOTP(c, repo) })
    auth.POST("/resend-otp", func(c *gin.Context){ handleResendOTP(c, repo) })
}

func handleRegister(c *gin.Context, repo *UserRepo) {
    var in RegisterInput
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()

    email := strings.ToLower(in.Email)
    hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to hash"}); return }
    userID, err := repo.CreateUser(ctx, in.Username, email, in.Phone, string(hash))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"email/username exists"}); return }

    code := generateOTP6()
    expires := time.Now().Add(10 * time.Minute)
    _, _ = repo.CreateOTP(ctx, userID, in.SendVia, code, expires)

    if in.SendVia == "wa" {
        go sendWhatsAppOTP(in.Phone, code)
    } else {
        go sendEmailOTP(email, code)
    }

    c.JSON(http.StatusCreated, gin.H{"user_id": userID, "message": "user created, OTP sent"})
}

func handleLogin(c *gin.Context, repo *UserRepo) {
    var in LoginInput
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()

    user, pass, err := repo.GetByEmail(ctx, strings.ToLower(in.Email))
    if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return }
    if bcrypt.CompareHashAndPassword([]byte(pass), []byte(in.Password)) != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return }

    // issue token with role
    token, _ := issueJWT(user.ID, user.Username, user.Email, user.Role)
    c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func handleVerifyOTP(c *gin.Context, repo *UserRepo) {
    var in VerifyOTPInput
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()

    otp, err := repo.GetValidOTP(ctx, in.UserID, in.Channel, in.Code)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid or expired code"}); return }
    _ = repo.MarkOTPUsed(ctx, otp.ID)
    _ = repo.SetVerified(ctx, otp.UserID)
    c.JSON(http.StatusOK, gin.H{"message":"verified"})
}

func handleResendOTP(c *gin.Context, repo *UserRepo) {
    var in struct {
        Email   string `json:"email"`
        UserID  string `json:"user_id"`
        Channel string `json:"channel"`
    }
    if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()

    var user User
    var err error
    if in.UserID != "" {
        user, err = repo.GetByID(ctx, in.UserID)
    } else {
        user, _, err = repo.GetByEmail(ctx, strings.ToLower(in.Email))
    }
    if err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"user not found"}); return }

    code := generateOTP6()
    expires := time.Now().Add(10 * time.Minute)
    _, _ = repo.CreateOTP(ctx, user.ID, in.Channel, code, expires)
    if in.Channel == "wa" && user.Phone != nil {
        go sendWhatsAppOTP(*user.Phone, code)
    } else {
        go sendEmailOTP(user.Email, code)
    }
    c.JSON(http.StatusOK, gin.H{"message":"otp resent"})
}
