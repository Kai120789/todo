package storage

import (
	"context"
	"fmt"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db *pgxpool.Pool
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
func (d *Storage) GetAllBoards() ([]models.Board, error) {
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

	return boards, nil
}

// функция получения доски
func (d *Storage) GetBoard(id uint) (*models.Board, error) {
	query := `SELECT * FROM boards WHERE name = $1 ORDER BY created_at`
	row := d.db.QueryRow(context.Background(), query, id)

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

// функция обновления доски
func (d *Storage) UpdateBoard(board dto.PostBoardDto) error {
	query := `UPDATE task SET name=$1, user_id=$2 WHERE id=$3`
	_, err := d.db.Exec(context.Background(), query, board.Name, board.User_id, board.ID)
	if err != nil {
		return err
	}

	return nil
}

// функция удаления доски
func (d *Storage) DeleteBoard(id string) error {
	query := `DELETE FROM board WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
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
func (d *Storage) GetTask(id uint) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1 ORDER BY updated_at`
	row := d.db.QueryRow(context.Background(), query, id)

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
func (d *Storage) GetAllTasks() ([]models.Task, error) {
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
	return tasks, nil
}

// функция обновления задачи
func (d *Storage) UpdateTask(task dto.PostTaskDto) error {
	query := `UPDATE task SET title=$1, description=$2, board_id=$3, status_id=$4, user_id=$5 WHERE id=$6`
	_, err := d.db.Exec(context.Background(), query, task.Title, task.Description, task.Board_id, task.Status_id, task.User_id, task.ID)
	if err != nil {
		return err
	}

	return nil
}

// функция удаления задачи
func (d *Storage) DeleteTask(id string) error {
	query := `DELETE FROM tasks WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// функция создания статуса
func (d *Storage) SetStatus() error {
	query := `INSERT INTO statuses (id, type) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

// функция удаления статуса
func (d *Storage) DeleteStatus(id string) error {
	query := `DELETE FROM statuses WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// функция для регистрации пользователя

// проверка, зарегистрирован ли пользователь

// удаление пользователя

func (d *Storage) Close() error {
	if d.db == nil {
		return nil
	}
	d.db.Close()
	return nil
}
