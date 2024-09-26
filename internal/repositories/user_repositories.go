package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"todo/internal/db"
	"todo/internal/models"

	_ "github.com/lib/pq" // Импорт PostgreSQL драйвера
)

// TaskRepository — структура репозитория для работы с задачами и пользователями
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository — конструктор для UserRepository
func NewUserRepository() *UserRepository {
	// Подключение к базе данных (измените настройки подключения на ваши)
	db, err := sql.Open("postgres", "host=localhost port=5432 user=youruser password=yourpassword dbname=taskdb sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	return &UserRepository{db: db}
}

func CreateUser(user models.User) (int, error) {
	var id int
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	err := db.GetDBPool().QueryRow(context.Background(), query, user.Username, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %v", err)
	}
	return id, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	err := db.GetDBPool().QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		return user, fmt.Errorf("error fetching user: %w", err)
	}
	return user, nil
}
