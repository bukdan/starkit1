package middleware

import (
	"context"
	"net/http"
	"strings"

	"gateway/utils"
)

type contextKey string

const UserCtxKey = contextKey("user")

// AuthMiddleware untuk GraphQL
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// inject ke context
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Ambil user dari context
func ForContext(ctx context.Context) *utils.Claims {
	raw, ok := ctx.Value(UserCtxKey).(*utils.Claims)
	if !ok {
		return nil
	}
	return raw
}
