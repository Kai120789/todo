package tokens

import (
	"errors"
	"time"
	"todo/internal/todo/config"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint, expiresAt time.Time) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("Error fetching config", zap.Error(err))
		return "", err
	}

	// Проверьте, что SecretKey — это строка
	if cfg.SecretKey == "" {
		zap.S().Error("Secret key is empty")
		return "", errors.New("secret key is empty")
	}

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(), // Здесь устанавливается время истечения
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.SecretKey)) // Приведение ключа к []byte
	if err != nil {
		zap.S().Errorf("Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
}
