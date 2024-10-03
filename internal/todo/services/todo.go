package services

import (
	"fmt"
	"strconv"
	"todo/internal/todo/api/tg"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

type TodoService struct {
	storage Storager
}

type Storager interface {
	SetBoard(body dto.PostBoardDto) (*models.Board, error)
	GetAllBoards() ([]models.Board, error)
	GetBoard(id uint) (*models.Board, error)
	UpdateBoard(body dto.PostBoardDto, id uint) (*models.Board, error)
	DeleteBoard(id uint) error

	User2Board(body dto.PostUser2BoardDto) error

	SetTask(body dto.PostTaskDto) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto, id uint) (*models.Task, error)
	DeleteTask(id uint) error

	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id uint) error

	RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error)
	AuthorizateUser(body dto.PostUserDto) (*models.UserToken, *uint, error)
	GetAuthUser(id uint) (*models.UserToken, error)
	UserLogout(id uint) error
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

	_, err := t.storage.SetBoard(body)
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

func (t *TodoService) UpdateBoard(body dto.PostBoardDto, id uint) error {
	if body.Name == "" {
		return fmt.Errorf("board Name is required")
	}

	_, err := t.storage.UpdateBoard(body, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) DeleteBoard(id string) error {
	if id == "" {
		return fmt.Errorf("board ID is required")
	}

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

func (t *TodoService) User2Board(body dto.PostUser2BoardDto) error {
	if body.UserId == "" {
		return fmt.Errorf("user ID is required")
	}

	if body.BoardId == "" {
		return fmt.Errorf("board ID is required")
	}

	err := t.storage.User2Board(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) SetTask(body dto.PostTaskDto) error {
	if body.Title == "" {
		return fmt.Errorf("task title cannot be empty")
	}

	task, err := t.storage.SetTask(body)
	if err != nil {
		return err
	}

	err = tg.Create(task)
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

func (t *TodoService) UpdateTask(body dto.PostTaskDto, id uint) error {
	if body.Title == "" {
		return fmt.Errorf("task Title is required")
	}

	_, err := t.storage.UpdateTask(body, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("task ID is required")
	}

	Uintid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	err = t.storage.DeleteTask(uint(Uintid))
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) SetStatus(body dto.PostStatusDto) error {
	err := t.storage.SetStatus(body)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) DeleteStatus(id string) error {
	if id == "" {
		return fmt.Errorf("status ID is required")
	}

	Uintid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	err = t.storage.DeleteStatus(uint(Uintid))
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) RegisterNewUser(body dto.PostUserDto) (*models.UserToken, error) {
	if body.Username == "" {
		return nil, fmt.Errorf("task title cannot be empty")
	}

	token, err := t.storage.RegisterNewUser(body)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TodoService) AuthorizateUser(body dto.PostUserDto) (*models.UserToken, *uint, error) {
	if body.Username == "" {
		return nil, nil, fmt.Errorf("task title cannot be empty")
	}

	token, id, err := t.storage.AuthorizateUser(body)
	if err != nil {
		return nil, nil, err
	}

	return token, id, nil
}

func (t *TodoService) GetAuthUser(id uint) (*models.UserToken, error) {
	token, err := t.storage.GetAuthUser(uint(id))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TodoService) UserLogout(id uint) error {
	err := t.storage.UserLogout(uint(id))
	if err != nil {
		return err
	}

	return nil
}
