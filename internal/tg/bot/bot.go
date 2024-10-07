package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type TgService struct {
	tgstorage Storager
}

type Storager interface {
	RegisterUser(upd int, tgName string, chatID int64) error
}

func New(stor Storager, logger *zap.Logger) *TgService {
	return &TgService{
		tgstorage: stor,
	}
}

func (tg *TgService) CheckUpdates(bot *tgbotapi.BotAPI) {

}
