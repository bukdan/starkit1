package middleware

import (
	"net/http"
	"strings"
	"user-service/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware -> validasi JWT token dari request
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// format token: "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		// ambil token string
		tokenStr := tokenParts[1]

		// parse token pakai helper utils
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// simpan user_id dari token ke context
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
