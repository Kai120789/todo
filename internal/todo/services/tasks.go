package services

import (
	"strconv"
	"todo/internal/todo/api"
	"todo/internal/todo/dto"
	"todo/internal/todo/models"

	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
)

type TasksService struct {
	storage TasksStorager
}

type TasksStorager interface {
	SetTask(body dto.PostTaskDto) (*models.Task, error)
	GetTask(id uint) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(body dto.PostTaskDto, id uint) (*models.Task, error)
	DeleteTask(id uint) error
	GetChatID(task *models.Task) (*int64, error)
	GetMyTasks(tgName string, status int) ([]models.Task, *int64, error)
	ChangeEndedTasksStatus() error
	GetAllUsers() ([]models.TgUser, error)
}

func NewTasksService(stor TasksStorager, logger *zap.Logger) *TasksService {
	return &TasksService{
		storage: stor,
	}
}

func (t *TasksService) SetTask(body dto.PostTaskDto) error {
	task, err := t.storage.SetTask(body)
	if err != nil {
		return err
	}

	chatID, err := t.storage.GetChatID(task)
	if err != nil {
		return err
	}

	err = api.Create(*task, *chatID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TasksService) GetTask(id uint) (*models.Task, error) {
	task, err := t.storage.GetTask(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *TasksService) GetAllTasks() ([]models.Task, error) {
	tasks, err := t.storage.GetAllTasks()
	if err != nil {
		return nil, err
	}

	return tasks, err
}

func (t *TasksService) UpdateTask(body dto.PostTaskDto, id uint) error {
	_, err := t.storage.UpdateTask(body, id)
	if err != nil {
		return err
	}

	return nil
}

func (t *TasksService) DeleteTask(id string) error {
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

func (t *TasksService) SendAllTasks(tgName string, chatID int64) error {
	message, _, err := t.storage.GetMyTasks(tgName, 1)
	if err != nil {
		zap.S().Error("Ошибка получения задач для пользователя", zap.String("tgName", tgName), zap.Error(err))
		return err
	}
	err = api.SendDailyReports(message, chatID, 1)
	if err != nil {
		return err
	}

	message, _, err = t.storage.GetMyTasks(tgName, 2)
	if err != nil {
		zap.S().Error("Ошибка получения выполненных задач для пользователя", zap.String("tgName", tgName), zap.Error(err))
		return err
	}
	err = api.SendDailyReports(message, chatID, 2)
	if err != nil {
		return err
	}

	return nil
}

func (t *TasksService) SendDailyReport() {
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

func (t *TasksService) StartScheduler() {
	gocron.Every(1).Day().At("00:00").Do(func() {
		t.SendDailyReport()
	})

	go func() {
		<-gocron.Start()
	}()
}
