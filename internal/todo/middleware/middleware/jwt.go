package middleware

import (
	"context"
	"net/http"
	"todo/internal/todo/config"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// Claims для JWT
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// Middleware для проверки Access токена
func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.GetConfig()
		if err != nil {
			zap.S().Fatalf("get config error", zap.Error(err))
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return cfg.SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
