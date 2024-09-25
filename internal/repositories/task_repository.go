package repositories

import (
	"database/sql"
	"log"
	"todo/internal/models"

	_ "github.com/lib/pq" // Используем PostgreSQL драйвер
)

type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository — конструктор для TaskRepository
func NewTaskRepository() *TaskRepository {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=youruser password=yourpassword dbname=taskdb sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	return &TaskRepository{db: db}
}

// GetTasksByUserID — получение всех задач для конкретного пользователя
func (repo *TaskRepository) GetTasksByUserID(userID int) ([]models.Task, error) {
	rows, err := repo.db.Query("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.StatusID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID — получение задачи по ID и userID
func (repo *TaskRepository) GetTaskByID(userID, taskID int) (*models.Task, error) {
	row := repo.db.QueryRow("SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)

	var task models.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.StatusID, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}
