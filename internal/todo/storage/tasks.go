package storage

import (
	"context"
	"fmt"
	"strconv"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TasksStorage struct {
	db *pgxpool.Pool
}

type TasksStorager interface {
	SetTask(body dto.PostTaskDto) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto) (*models.Task, error)
	DeleteTask(id uint) error
	GetChatID(task *models.Task) (*int64, error)
	GetMyTasks(tgName string, status int) ([]models.Task, *int64, error)
	ChangeEndedTasksStatus() error
	GetAllUsers() ([]models.TgUser, error)
}

func NewTasksStore(Conn *pgxpool.Pool, log *zap.Logger) *TasksStorage {
	return &TasksStorage{db: Conn}
}

// set task
func (d *TasksStorage) SetTask(body dto.PostTaskDto) (*models.Task, error) {
	userId, err := strconv.ParseUint(body.UserId, 10, 32)
	if err != nil {
		return nil, err
	}

	boardId, err := strconv.ParseUint(body.BoardId, 10, 32)
	if err != nil {
		return nil, err
	}

	var id uint
	query := `INSERT INTO tasks (title, description, board_id, status_id, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = d.db.QueryRow(context.Background(), query, body.Title, body.Description, boardId, 1, userId).Scan(&id)
	if err != nil {
		return nil, err
	}

	taskRet, err := d.GetTask(uint(id))
	if err != nil {
		return nil, err
	}

	return taskRet, nil
}

// get task
func (d *TasksStorage) GetTask(id uint) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1`
	row := d.db.QueryRow(context.Background(), query, id)

	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.BoardId, &task.StatusId, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

// get all tasks
func (d *TasksStorage) GetAllTasks() ([]models.Task, error) {
	query := `SELECT id, title, description, board_id, status_id, user_id, created_at, updated_at FROM tasks ORDER BY updated_at`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.BoardId, &task.StatusId, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// update task
func (d *TasksStorage) UpdateTask(body dto.PostTaskDto, id uint) (*models.Task, error) {
	userId, err := strconv.ParseUint(body.UserId, 10, 32)
	if err != nil {
		return nil, err
	}

	boardId, err := strconv.ParseUint(body.BoardId, 10, 32)
	if err != nil {
		return nil, err
	}

	query := `UPDATE tasks SET title=$1, description=$2, board_id=$3, status_id=$4, user_id=$5, updated_at=NOW() WHERE id=$6`
	_, err = d.db.Exec(context.Background(), query, body.Title, body.Description, boardId, body.StatusId, userId, id)
	if err != nil {
		return nil, err
	}

	taskRet, err := d.GetTask(uint(id))
	if err != nil {
		return nil, err
	}

	return taskRet, nil
}

// delete task
func (d *TasksStorage) DeleteTask(id uint) error {
	query := `DELETE FROM tasks WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *TasksStorage) GetChatID(task *models.Task) (*int64, error) {
	var chatID int64
	query := `SELECT chat_id FROM users WHERE id=$1`
	err := d.db.QueryRow(context.Background(), query, task.UserId).Scan(&chatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("username not found")
		}
		return nil, err
	}

	return &chatID, err
}

func (d *TasksStorage) GetMyTasks(tgName string, status int) ([]models.Task, *int64, error) {
	var id, chatID uint

	query := `SELECT id, chat_id FROM users WHERE tg_name=$1`
	err := d.db.QueryRow(context.Background(), query, tgName).Scan(&id, &chatID)
	fmt.Println("chatID from DB:", chatID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	query = `SELECT * FROM tasks WHERE user_id=$1 and status_id=$2 ORDER BY updated_at`
	rows, err := d.db.Query(context.Background(), query, id, status)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.BoardId, &task.StatusId, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, nil, err
		}
		tasks = append(tasks, task)
	}

	intChatID := int64(chatID)

	return tasks, &intChatID, nil
}

func (d *TasksStorage) ChangeEndedTasksStatus() error {
	query := `UPDATE tasks SET status_id = 3 WHERE status_id = 2`
	_, err := d.db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса задач: %w", err)
	}

	return nil
}

func (d *TasksStorage) GetAllUsers() ([]models.TgUser, error) {
	query := `SELECT id, tg_name, chat_id FROM users`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.TgUser
	for rows.Next() {
		var user models.TgUser
		if err := rows.Scan(&user.ID, &user.TgName, &user.ChatID); err != nil {
			fmt.Println("Error during scanning row:", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
