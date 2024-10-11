package services

import (
	"strconv"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

type BoardsService struct {
	storage BoardsStorager
}

type BoardsStorager interface {
	SetBoard(body dto.PostBoardDto) (*models.Board, error)
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto, id uint) (*models.Board, error)
	DeleteBoard(id uint) error
	User2Board(body dto.PostUser2BoardDto) error
}

func NewBoardsService(stor BoardsStorager, logger *zap.Logger) *BoardsService {
	return &BoardsService{
		storage: stor,
	}
}

func (t *BoardsService) SetBoard(body dto.PostBoardDto) (*models.Board, error) {
	boardRet, err := t.storage.SetBoard(body)
	if err != nil {
		return nil, err
	}

	return boardRet, nil
}

func (t *BoardsService) GetAllBoards() ([]models.Board, error) {
	boards, err := t.storage.GetAllBoards()
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (t *BoardsService) GetBoard(id uint) (*models.Board, error) {
	board, err := t.storage.GetBoard(id)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (t *BoardsService) UpdateBoard(body dto.PostBoardDto, id uint) error {
	_, err := t.storage.UpdateBoard(body, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *BoardsService) DeleteBoard(id string) error {
	Uintid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	err = t.storage.DeleteBoard(uint(Uintid))
	if err != nil {
		return err
	}

	return nil
}

func (t *BoardsService) User2Board(body dto.PostUser2BoardDto) error {
	err := t.storage.User2Board(body)
	if err != nil {
		return err
	}

	return nil
}
