package storage

import (
	"context"
	"fmt"
	"todo/internal/todo/dto"
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
func (d *Storage) SetBoard(board dto.PostBoardDto) error {
	query := `INSERT INTO boards (id, name, user_id) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(context.Background(), query, board.ID, board.Name, board.User_id)
	if err != nil {
		return err
	}

	return nil
}

// функция получения всех досок
func (d *Storage) GetAllBoards() (*[]models.Board, error) {
	query := `SELECT id, name, user_id, created_at from boards ORDER BY created_at`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var boards []models.Board
	for rows.Next() {
		var board models.Board
		err := rows.Scan(&board.ID, &board.Name, &board.User_id, &board.Created_at)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return &boards, nil
}

// функция получения доски
func (d *Storage) GetBoard(name string) (*models.Board, error) {
	query := `SELECT * FROM boards WHERE name = $1 ORDER BY created_at`
	row := d.db.QueryRow(context.Background(), query, name)

	var board models.Board
	err := row.Scan(&board.ID, &board.Name, &board.User_id, &board.Created_at)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &board, nil
}

// функция создания задачи
func (d *Storage) SetTask(task dto.PostTaskDto) error {
	query := `INSERT INTO tasks (id, title, description, board_id, status_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := d.db.Exec(context.Background(), query, task.ID, task.Title, task.Description, task.Board_id, task.Status_id, task.User_id)
	if err != nil {
		return err
	}

	return nil
}

// функция получения задачи
func (d *Storage) GetTask(name string) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE name = $1 ORDER BY updated_at`
	row := d.db.QueryRow(context.Background(), query, name)

	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Board_id, &task.Status_id, &task.User_id, &task.Created_at, &task.Updated_at)
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
	query := `SELECT id, title, description, board_id, status_id, user_id, created_at, updated_at FROM tasks ORDER BY updated_at`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Board_id, &task.Status_id, &task.User_id, &task.Created_at, &task.Updated_at)
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
