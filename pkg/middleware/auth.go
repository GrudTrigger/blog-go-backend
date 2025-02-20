package middleware

import (
	"backend/blog/configs"
	"backend/blog/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func IsAuthed(next http.Handler, config *configs.Configs) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if authedHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authedHeader, "Bearer ")
		if token == "" {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		isValid, data := jwt.NewJWT(config.Secret).Parse(token)
		if !isValid || data == nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}