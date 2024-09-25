package repositories

import (
	"database/sql"
	"fmt"
	"log"
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

var userRepo = NewUserRepository()

// CreateUser — создание пользователя
func CreateUser(user models.User) (int, error) {
	var id int
	err := userRepo.db.QueryRow(`
		INSERT INTO users (username, password) 
		VALUES ($1, $2) 
		RETURNING id
	`, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %v", err)
	}
	return id, nil
}

// GetUserByUsername — получение пользователя по имени пользователя
func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	row := userRepo.db.QueryRow(`
		SELECT id, username, password 
		FROM users 
		WHERE username = $1
	`, username)

	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, fmt.Errorf("error fetching user: %v", err)
	}

	return user, nil
}
