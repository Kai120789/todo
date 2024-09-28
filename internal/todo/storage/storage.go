package storage

import (
	"context"
	"fmt"
	"todo/internal/todo/models"

	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(Conn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{db: Conn, logger: log}
}

func Connection(connectionStr string, logger *zap.Logger) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	log.Println("Successfully connect to db")

	return db, nil
}

func (d *Storage) Ping() error {
	fmt.Println("ping")
	if d.db == nil {
		return fmt.Errorf("database is not connected")
	}
	return d.db.Ping(context.Background())
}

// функция добавления доски

// функция получения доски

// функция создания задачи

// функция получения задачи
func (d *Storage) GetTask(name string) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE name = $1 ORDER BY updated_at`
	row := d.db.QueryRow(context.Background(), query, name)

	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Created_at, &task.Updated_at)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

// функция получения всех задач
func (d *Storage) GetAllTasks() (*[]models.Task, error) {
	query := `SELECT id, title, description, created_at, updated_at FROM tasks ORDER BY updated_at`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Created_at, &task.Updated_at)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil
}

func (d *Storage) Close() error {
	if d.db == nil {
		return nil
	}
	d.db.Close()
	return nil
}
