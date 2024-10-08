package services

import (
	"fmt"
	"strconv"
	"todo/internal/todo/api"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/jasonlvhit/gocron"

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
	WriteRefreshToken(userId uint, refreshTokenValue string) error
	GetAuthUser(id uint) (*models.UserToken, error)
	UserLogout(id uint) error

	GetAllUsers() ([]models.TgUser, error)
	GetChatID(task *models.Task) (*int64, error)
	AddChatID(tgName string, chatID int64) error
	GetMyTasks(tgName string, status int) ([]models.Task, *int64, error)
	ChangeEndedTasksStatus() error
}

func New(stor Storager, logger *zap.Logger) *TodoService {
	return &TodoService{
		storage: stor,
	}
}

func (t *TodoService) SetBoard(body dto.PostBoardDto) (*models.Board, error) {
	if body.Name == "" {
		return nil, fmt.Errorf("board name cannot be empty")
	}

	boardRet, err := t.storage.SetBoard(body)
	if err != nil {
		return nil, err
	}

	return boardRet, nil
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

	chatID, err := t.storage.GetChatID(task)
	if err != nil {
		return err
	}

	err = api.Create(task, *chatID)
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
		return nil, nil, fmt.Errorf("username cannot be empty")
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

func (t *TodoService) WriteRefreshToken(userId uint, refreshTokenValue string) error {
	err := t.storage.WriteRefreshToken(userId, refreshTokenValue)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) AddChatID(tgName string, chatID int64) error {
	err := t.storage.AddChatID(tgName, chatID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoService) SendDailyReport() {
	users, err := t.storage.GetAllUsers()
	if err != nil {
		zap.L().Error("Ошибка при получении пользователей", zap.Error(err))
		return
	}

	for _, user := range users {
		message, _, err := t.storage.GetMyTasks(user.TgName, 1)
		if err != nil {
			zap.S().Error("Ошибка получения задач для пользователя", zap.String("tgName", user.TgName), zap.Error(err))
			continue
		}
		api.SendDailyReports(message, user.ChatID, 1)

		message, _, err = t.storage.GetMyTasks(user.TgName, 2)
		if err != nil {
			zap.S().Error("Ошибка получения выполненных задач для пользователя", zap.String("tgName", user.TgName), zap.Error(err))
			continue
		}
		api.SendDailyReports(message, user.ChatID, 2)

		err = t.storage.ChangeEndedTasksStatus()
		if err != nil {
			zap.L().Error("Ошибка обновления статуса задач", zap.String("tgName", user.TgName), zap.Error(err))
		}
	}
}

func (t *TodoService) StartScheduler() {
	gocron.Every(1).Day().At("00:00").Do(func() {
		t.SendDailyReport()
	})

	go func() {
		<-gocron.Start()
	}()
}
