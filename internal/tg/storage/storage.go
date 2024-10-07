package tgstorage

import (
	"context"
	"fmt"
	"todo/internal/todo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db *pgxpool.Pool
}

type Storager interface {
	RegisterUser(upd int, tgName string, chatID int64) error
	SendTask()
	DailyReport()
}

func New(Conn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{db: Conn}
}

func Connection(connectionStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}

func (d *Storage) RegisterUser(upd int, tgName string, chatID int64) error {
	query := `INSERT INTO tg_id (tg_name, chat_id) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query, tgName, chatID)
	if err != nil {
		fmt.Printf("Error registering user: %v\n", err)
		return err
	}

	return nil
}

func (d *Storage) GetMyTasks(tgName string) ([]models.Task, int64, error) {
	var id uint

	query := `SELECT id FROM users WHERE tg_name=$1`
	err := d.db.QueryRow(context.Background(), query, tgName).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, 0, fmt.Errorf("user not found")
		}
		return nil, 0, err
	}

	query = `SELECT id, title, description, board_id, status_id, created_at, updated_at FROM tasks WHERE user_id=$1 ORDER BY updated_at`
	rows, err := d.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.BoardId, &task.StatusId, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	var chatId int64
	query = `SELECT chat_id FROM tg_id WHERE tg_name=$1`
	err = d.db.QueryRow(context.Background(), query, tgName).Scan(&chatId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, 0, fmt.Errorf("chat not found")
		}
		return nil, 0, err
	}

	return tasks, chatId, nil
}

func CreateTask() {

}

func DailyReport() {

}
