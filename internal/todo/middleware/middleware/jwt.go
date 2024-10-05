package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"todo/internal/todo/config"
	"todo/pkg/logger"

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
		// Получаем конфигурацию
		cfg, err := config.GetConfig()
		if err != nil {
			zap.S().Fatalf("Ошибка получения конфигурации", zap.Error(err))
			return
		}

		zapLog, err := logger.New(cfg.LogLevel)
		if err != nil {
			zap.S().Fatalf("init logger error", zap.Error(err))
		}

		log := zapLog.ZapLogger

		// Извлекаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
			return
		}

		// Проверяем, что токен начинается с 'Bearer '
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		// Убираем префикс 'Bearer '
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Выводим полученный токен
		log.Info("Полученный токен", zap.String("token", tokenString))

		// Определяем, куда будет сохраняться информация из токена
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Убедитесь, что метод подписи верный
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return []byte(cfg.SecretKey), nil
		})

		// Проверяем валидность токена
		if err != nil {
			log.Error("Ошибка разбора токена: %v", zap.Error(err)) // Вывод ошибки разбора токена
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Error("Токен недействителен") // Сообщение о том, что токен недействителен
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Выводим разобранные claims
		zap.S().Infof("Разобранные claims: %+v", claims)

		// Добавляем user_id в контекст запроса
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)

		// Передаём запрос дальше по цепочке
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
