package utils

import (
	"fmt"
	"todo/internal/tg/dto"
)

func FormatTasksMessage(tasks []dto.MessDto) (*int64, string, string) {
	if len(tasks) == 0 {
		return nil, "У вас нет задач.", "У вас нет задачб завершенных сегодня."
	}

	var chatID *int64

	message := "Ваши задачи:\n\n"
	messageEnded := "Ваши завершенные задачи:\n\n"
	for i, task := range tasks {
		if task.StatusId == 1 {
			message += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %s\n\n", i+1, task.Title, task.Description, "в процессе")
		}

		chatID = &task.ChatId
	}

	for i, task := range tasks {
		if task.StatusId == 2 {
			messageEnded += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %s\n\n", i+1, task.Title, task.Description, "выполнено")
		}

		chatID = &task.ChatId
	}

	return chatID, message, messageEnded
}
