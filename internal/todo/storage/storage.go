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
func (d *Storage) SetBoard(board dto.PostBoardDto) (*models.Board, error) {
	id, err := strconv.ParseUint(board.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	userId, err := strconv.ParseUint(board.UserId, 10, 32)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO boards (id, name, user_id) VALUES ($1, $2, $3)`
	_, err = d.db.Exec(context.Background(), query, id, board.Name, userId)
	if err != nil {
		return nil, err
	}

	// заполнение таблицы boards_users для связи многие ко многим

	boardRet, err := d.GetBoard(uint(id))
	if err != nil {
		return nil, err
	}

	return boardRet, nil
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
		err := rows.Scan(&board.ID, &board.Name, &board.UserId, &board.CreatedAt)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}

// функция получения доски
func (d *Storage) GetBoard(id uint) (*models.Board, error) {
	query := `SELECT * FROM boards WHERE id = $1 ORDER BY created_at`
	row := d.db.QueryRow(context.Background(), query, id)

	var board models.Board
	err := row.Scan(&board.ID, &board.Name, &board.UserId, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &board, nil
}

// функция обновления доски
func (d *Storage) UpdateBoard(board dto.PostBoardDto) (*models.Board, error) {
	query := `UPDATE boards SET name=$1, user_id=$2, updated_at=NOW() WHERE id=$3`
	_, err := d.db.Exec(context.Background(), query, board.Name, board.UserId, board.ID)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(board.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	boardRet, err := d.GetBoard(uint(id))
	if err != nil {
		return nil, err
	}

	return boardRet, nil
}

// функция удаления доски
func (d *Storage) DeleteBoard(id uint) error {
	query := `DELETE FROM boards WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// функция создания задачи
func (d *Storage) SetTask(task dto.PostTaskDto) (*models.Task, error) {
	id, err := strconv.ParseUint(task.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	userId, err := strconv.ParseUint(task.UserId, 10, 32)
	if err != nil {
		return nil, err
	}

	boardId, err := strconv.ParseUint(task.BoardId, 10, 32)
	if err != nil {
		return nil, err
	}

	statusId, err := strconv.ParseUint(task.StatusId, 10, 32)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO tasks (id, title, description, board_id, status_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = d.db.Exec(context.Background(), query, id, task.Title, task.Description, boardId, statusId, userId)
	if err != nil {
		return nil, err
	}

	taskRet, err := d.GetTask(uint(id))
	if err != nil {
		return nil, err
	}

	return taskRet, nil
}

// функция получения задачи
func (d *Storage) GetTask(id uint) (*models.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1 ORDER BY updated_at`
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
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.BoardId, &task.StatusId, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// функция обновления задачи
func (d *Storage) UpdateTask(task dto.PostTaskDto) (*models.Task, error) {
	query := `UPDATE tasks SET title=$1, description=$2, board_id=$3, status_id=$4, user_id=$5, updated_at=NOW() WHERE id=$6`
	_, err := d.db.Exec(context.Background(), query, task.Title, task.Description, task.BoardId, task.StatusId, task.UserId, task.ID)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(task.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	taskRet, err := d.GetTask(uint(id))
	if err != nil {
		return nil, err
	}

	return taskRet, nil
}

// функция удаления задачи
func (d *Storage) DeleteTask(id uint) error {
	query := `DELETE FROM tasks WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// функция создания статуса
func (d *Storage) SetStatus(status dto.PostStatusDto) error {
	query := `INSERT INTO statuses (id, type) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query, status.ID, status.Type)
	if err != nil {
		return err
	}

	return nil
}

// функция удаления статуса
func (d *Storage) DeleteStatus(id uint) error {
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
