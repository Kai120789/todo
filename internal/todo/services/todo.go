package services

import (
	"fmt"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

type TodoService struct {
	storage Storager
}

type Storager interface {
	SetBoard(body dto.PostBoardDto) error
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto) error
	DeleteBoard(id string) error
	SetTask(body dto.PostTaskDto) error
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto) error
	DeleteTask(id string) error
	SetStatus() error
	DeleteStatus(id string) error
}

func New(stor Storager, logger *zap.Logger) *TodoService {
	return &TodoService{
		storage: stor,
	}
}

func (t *TodoService) SetBoard(body dto.PostBoardDto) error {
	if body.Name == "" {
		return fmt.Errorf("board name cannot be empty")
	}

	err := t.storage.SetBoard(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) GetAllBoards() ([]models.Board, error) {
	boards, err := t.storage.GetAllBoards()
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (t *TodoService) GetBoard(id uint) (*models.Board, error) {
	board, err := t.storage.GetBoard(id)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (t *TodoService) UpdateBoard(body dto.PostBoardDto) error {
	if body.ID == "" {
		return fmt.Errorf("board ID is required")
	}

	err := t.storage.UpdateBoard(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) DeleteBoard(id string) error {
	err := t.storage.DeleteBoard(id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) SetTask(body dto.PostTaskDto) error {
	if body.Title == "" {
		return fmt.Errorf("task title cannot be empty")
	}

	err := t.storage.SetTask(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) GetTask(id uint) (*models.Task, error) {
	task, err := t.storage.GetTask(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *TodoService) GetAllTasks() ([]models.Task, error) {
	tasks, err := t.storage.GetAllTasks()
	if err != nil {
		return nil, err
	}

	return tasks, err
}

func (t *TodoService) UpdateTask(body dto.PostTaskDto) error {
	if body.ID == "" {
		return fmt.Errorf("task ID is required")
	}

	err := t.storage.UpdateTask(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("task ID is required")
	}
	err := t.storage.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) SetStatus() error {
	err := t.storage.SetStatus()
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) DeleteStatus(id string) error {
	if id == "" {
		return fmt.Errorf("status ID is required")
	}
	err := t.storage.DeleteStatus(id)
	if err != nil {
		return err
	}

	return nil
}
