package repositories

import (
	"database/sql"
	"todo/internal/models"
)

func CreateBoard(db *sql.DB, board models.Board) (int, error) {
	var id int
	query := `INSERT INTO boards (name, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(query, board.Name, board.OwnerID, board.CreatedAt, board.UpdatedAt).Scan(&id)
	return id, err
}

func GetBoards(db *sql.DB, userID int) ([]models.Board, error) {
	var boards []models.Board
	query := `SELECT id, name, owner_id, created_at, updated_at FROM boards WHERE owner_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var board models.Board
		err := rows.Scan(&board.ID, &board.Name, &board.OwnerID, &board.CreatedAt, &board.UpdatedAt)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func GetBoardByID(db *sql.DB, id int, userID int) (*models.Board, error) {
	var board models.Board
	query := `SELECT id, name, owner_id, created_at, updated_at FROM boards WHERE id = $1 AND owner_id = $2`
	err := db.QueryRow(query, id, userID).Scan(&board.ID, &board.Name, &board.OwnerID, &board.CreatedAt, &board.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &board, err
}

func UpdateBoard(db *sql.DB, board models.Board) error {
	query := `UPDATE boards SET name = $1, updated_at = $2 WHERE id = $3 AND owner_id = $4`
	_, err := db.Exec(query, board.Name, board.UpdatedAt, board.ID, board.OwnerID)
	return err
}

func DeleteBoard(db *sql.DB, id int, userID int) error {
	query := `DELETE FROM boards WHERE id = $1 AND owner_id = $2`
	_, err := db.Exec(query, id, userID)
	return err
}
