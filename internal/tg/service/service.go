package tgservice

import (
	"fmt"
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
}

func New(stor Storager, logger *zap.Logger) *TgService {
	return &TgService{
		tgstorage: stor,
	}
}

func (tg *TgService) GetMyTasks(tgName string) (string, int64, error) {
	tasks, chatId, err := tg.tgstorage.GetMyTasks(tgName)
	if err != nil {
		return "", 0, fmt.Errorf("err: %s", err)
	}

	res := utils.FormatTasksMessage(tasks)

	return res, chatId, nil
}
