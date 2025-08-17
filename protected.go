
package main

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProtectedRoutes(r *gin.Engine, db *pgxpool.Pool) {
    repo := NewUserRepo(db)
    v := r.Group("/api/v1")
    v.GET("/me", AuthMiddleware(), func(c *gin.Context) {
        uid := c.GetString("sub")
        ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()
        user, err := repo.GetByID(ctx, uid)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed"}) ; return }
        c.JSON(http.StatusOK, user)
    })
    v.PUT("/me", AuthMiddleware(), func(c *gin.Context) {
        uid := c.GetString("sub")
        var in struct{ Username, Phone, AvatarURL string }
        if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
        ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()
        _ = repo.UpdateProfile(ctx, uid, in.Username, in.Phone, in.AvatarURL)
        c.JSON(http.StatusOK, gin.H{"message":"updated"})
    })

    admin := v.Group("/admin", AuthMiddleware(), AdminOnly())
    admin.GET("/users", func(c *gin.Context) {
        limit, _ := strconv.Atoi(c.Query("limit")); if limit==0 { limit=50 }
        offset, _ := strconv.Atoi(c.Query("offset"))
        ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()
        users, _ := repo.ListUsers(ctx, limit, offset)
        c.JSON(http.StatusOK, users)
    })
    admin.GET("/users/:id", func(c *gin.Context){
        ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second); defer cancel()
        u, err := repo.GetByID(ctx, c.Param("id"))
        if err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
        c.JSON(http.StatusOK, u)
    })
    admin.PUT("/users/:id", func(c *gin.Context){
        var in struct{ Username, Email, Phone, Role string; IsVerified bool }
        if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()}); return }
        _ = repo.UpdateUserAdmin(c.Request.Context(), c.Param("id"), in.Username, in.Email, in.Phone, in.Role, in.IsVerified)
        c.JSON(http.StatusOK, gin.H{"message":"updated"})
    })
    admin.DELETE("/users/:id", func(c *gin.Context){
        _ = repo.DeleteUser(c.Request.Context(), c.Param("id"))
        c.JSON(http.StatusOK, gin.H{"message":"deleted"})
    })
}
