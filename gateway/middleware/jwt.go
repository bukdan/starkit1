package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const AuthTokenKey contextKey = "authToken"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			// inject token ke context
			ctx := context.WithValue(r.Context(), AuthTokenKey, token)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
