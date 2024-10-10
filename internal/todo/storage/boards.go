package storage

import (
	"context"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BoardsStorage struct {
	db *pgxpool.Pool
}

type BoardsStorager interface {
	SetBoard(body dto.PostBoardDto) (*models.Board, error)
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto) (*models.Board, error)
	DeleteBoard(id uint) error
}

func NewBoardsStore(Conn *pgxpool.Pool, log *zap.Logger) *BoardsStorage {
	return &BoardsStorage{db: Conn}
}

// add board
func (d *BoardsStorage) SetBoard(body dto.PostBoardDto) (*models.Board, error) {
	query := `INSERT INTO boards (name) VALUES ($1) RETURNING id`

	var id uint
	err := d.db.QueryRow(context.Background(), query, body.Name).Scan(&id)
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
func (d *BoardsStorage) GetAllBoards() ([]models.Board, error) {
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
func (d *BoardsStorage) GetBoard(id uint) (*models.Board, error) {
	query := `SELECT * FROM boards WHERE id = $1`
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
func (d *BoardsStorage) UpdateBoard(body dto.PostBoardDto, id uint) (*models.Board, error) {
	query := `UPDATE boards SET name=$1, updated_at=NOW() WHERE id=$2`
	_, err := d.db.Exec(context.Background(), query, body.Name, id)
	if err != nil {
		return nil, err
	}

	boardRet, err := d.GetBoard(id)
	if err != nil {
		return nil, err
	}

	return boardRet, nil
}

// delete board
func (d *BoardsStorage) DeleteBoard(id uint) error {
	query := `DELETE FROM boards WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

// add user to board
func (d *BoardsStorage) User2Board(body dto.PostUser2BoardDto) error {
	query := `INSERT INTO boards_users (user_id, board_id) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query, body.UserId, body.BoardId)
	if err != nil {
		return err
	}

	return nil

}
