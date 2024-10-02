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

type Storager interface {
	SetBoard(body dto.PostBoardDto) (*models.Board, error)
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto) (*models.Board, error)
	DeleteBoard(id uint) error

	User2Board(body dto.PostUser2BoardDto) error

	SetTask(body dto.PostTaskDto) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto) (*models.Task, error)
	DeleteTask(id uint) error

	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id uint) error
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

// add board
func (d *Storage) SetBoard(body dto.PostBoardDto) (*models.Board, error) {
	id, err := strconv.ParseUint(body.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	var maxID int
	err = d.db.QueryRow(context.Background(), "SELECT COALESCE(MAX(id), 0) FROM boards_users").Scan(&maxID)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO boards_users (id, user_id, board_id) VALUES ($1, $2, $3)`
	_, err = d.db.Exec(context.Background(), query, maxID+1, body.UserId, id)
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO boards (id, name) VALUES ($1, $2)`
	_, err = d.db.Exec(context.Background(), query, id, body.Name)
	if err != nil {
		return nil, err
	}

	boardRet, err := d.GetBoard(uint(id))
	if err != nil {
		return nil, err
	}

	return boardRet, nil
}

// get all boards
func (d *Storage) GetAllBoards() ([]models.Board, error) {
	query := `SELECT id, name, created_at, updated_at from boards ORDER BY created_at`
	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var boards []models.Board
	for rows.Next() {
		var board models.Board
		err := rows.Scan(&board.ID, &board.Name, &board.CreatedAt, &board.UpdatedAt)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}

// get board
func (d *Storage) GetBoard(id uint) (*models.Board, error) {
	query := `SELECT * FROM boards WHERE id = $1 ORDER BY created_at`
	row := d.db.QueryRow(context.Background(), query, id)

	var board models.Board
	err := row.Scan(&board.ID, &board.Name, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &board, nil
}

// update board
func (d *Storage) UpdateBoard(body dto.PostBoardDto) (*models.Board, error) {
	query := `UPDATE boards SET name=$1, updated_at=NOW() WHERE id=$2`
	_, err := d.db.Exec(context.Background(), query, body.Name, body.ID)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(body.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	boardRet, err := d.GetBoard(uint(id))
	if err != nil {
		return nil, err
	}

	return boardRet, nil
}

// delete board
func (d *Storage) DeleteBoard(id uint) error {
	query := `DELETE FROM boards WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// add user to board
func (d *Storage) User2Board(body dto.PostUser2BoardDto) error {
	query := `INSERT INTO boards_users (id, user_id, board_id) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(context.Background(), query, body.ID, body.UserId, body.BoardId)
	if err != nil {
		return err
	}

	return nil

}

// set task
func (d *Storage) SetTask(body dto.PostTaskDto) (*models.Task, error) {
	id, err := strconv.ParseUint(body.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	userId, err := strconv.ParseUint(body.UserId, 10, 32)
	if err != nil {
		return nil, err
	}

	boardId, err := strconv.ParseUint(body.BoardId, 10, 32)
	if err != nil {
		return nil, err
	}

	statusId, err := strconv.ParseUint(body.StatusId, 10, 32)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO tasks (id, title, description, board_id, status_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = d.db.Exec(context.Background(), query, id, body.Title, body.Description, boardId, statusId, userId)
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

// get all tasks
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

// update task
func (d *Storage) UpdateTask(body dto.PostTaskDto) (*models.Task, error) {
	id, err := strconv.ParseUint(body.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	userId, err := strconv.ParseUint(body.UserId, 10, 32)
	if err != nil {
		return nil, err
	}

	boardId, err := strconv.ParseUint(body.BoardId, 10, 32)
	if err != nil {
		return nil, err
	}

	statusId, err := strconv.ParseUint(body.StatusId, 10, 32)
	if err != nil {
		return nil, err
	}

	query := `UPDATE tasks SET title=$1, description=$2, board_id=$3, status_id=$4, user_id=$5, updated_at=NOW() WHERE id=$6`
	_, err = d.db.Exec(context.Background(), query, body.Title, body.Description, boardId, statusId, userId, id)
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
func (d *Storage) DeleteTask(id uint) error {
	query := `DELETE FROM tasks WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// create status
func (d *Storage) SetStatus(body dto.PostStatusDto) error {
	query := `INSERT INTO statuses (id, type) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query, body.ID, body.Type)
	if err != nil {
		return err
	}

	return nil
}

// delete status
func (d *Storage) DeleteStatus(id uint) error {
	query := `DELETE FROM statuses WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// register new user
func (d *Storage) RegisterNewUser(body dto.PostUserDto) (*models.User, error) {
	query := `INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(context.Background(), query, body.ID, body.Username, body.PasswordHash)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(body.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	userRet, err := d.GetAuthUser(uint(id))
	if err != nil {
		return nil, err
	}

	return userRet, nil
}

// login user
func (d *Storage) AuthorizateUser(body dto.GetUserDto) (*models.User, error) {

}

// get auth user
func (d *Storage) GetAuthUser(id uint) (*models.User, error) {

}

// logout user
func (d *Storage) UserLogout(id uint) error {

}

func (d *Storage) Close() error {
	if d.db == nil {
		return nil
	}
	d.db.Close()
	return nil
}
