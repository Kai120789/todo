package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"todo/internal/models"
	"todo/internal/repositories"
	"todo/internal/utils"
)

var db *sql.DB

func InitBoardHandler(database *sql.DB) {
	db = database
}

func CreateBoard(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)
	var board models.Board

	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	board.OwnerID = userID
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	id, err := repositories.CreateBoard(db, board)
	if err != nil {
		http.Error(w, "Error creating board", http.StatusInternalServerError)
		return
	}

	board.ID = id
	json.NewEncoder(w).Encode(board)
}

func GetBoards(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)

	boards, err := repositories.GetBoards(db, userID)
	if err != nil {
		http.Error(w, "Error fetching boards", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(boards)
}

func GetBoardByID(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)
	boardID, err := utils.ParseIDFromURL(r)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	board, err := repositories.GetBoardByID(db, boardID, userID)
	if err != nil {
		http.Error(w, "Error fetching board", http.StatusInternalServerError)
		return
	}

	if board == nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(board)
}

func UpdateBoard(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)
	boardID, err := utils.ParseIDFromURL(r)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	var board models.Board
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	board.ID = boardID
	board.OwnerID = userID
	board.UpdatedAt = time.Now()

	err = repositories.UpdateBoard(db, board)
	if err != nil {
		http.Error(w, "Error updating board", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(board)
}

func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)
	boardID, err := utils.ParseIDFromURL(r)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	err = repositories.DeleteBoard(db, boardID, userID)
	if err != nil {
		http.Error(w, "Error deleting board", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
