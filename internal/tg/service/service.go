package tgservice

import (
	"todo/internal/tg/utils"
	"todo/internal/todo/models"

	"go.uber.org/zap"
)

type TgService struct {
	tgstorage Storager
}

type Storager interface {
	RegisterUser(upd int, tgName string, chatID int64) error
	GetMyTasks(tgName string) ([]models.Task, int64, error)
	GetMyEndedTasks(tgName string) ([]models.Task, int64, error)
	ChangeEndedTasksStatus() error
}

func New(stor Storager, logger *zap.Logger) *TgService {
	return &TgService{
		tgstorage: stor,
	}
}

func (tg *TgService) GetMyTasks(tgName string) (string, int64, error) {
	tasks, chatId, err := tg.tgstorage.GetMyTasks(tgName)
	if err != nil {
		return "", 0, err
	}

	res := utils.FormatTasksMessage(tasks)

	return res, chatId, nil
}

func (tg *TgService) GetMyEndedTasks(tgName string) (string, int64, error) {
	tasks, chatId, err := tg.tgstorage.GetMyEndedTasks(tgName)
	if err != nil {
		return "", 0, err
	}

	res := utils.FormatEndedTasksMessage(tasks)

	return res, chatId, nil
}

func (tg *TgService) ChangeEndedTasksStatus() error {
	err := tg.tgstorage.ChangeEndedTasksStatus()
	if err != nil {
		return err
	}

	return nil
}
