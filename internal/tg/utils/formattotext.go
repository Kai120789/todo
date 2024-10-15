package utils

import (
	"fmt"
	"todo/internal/tg/dto"
)

func FormatTasksMessage(tasks []dto.MessDto) (*int64, string) {
	if len(tasks) == 0 {
		return nil, "У вас нет задач"
	}

	chatID := tasks[0].ChatId

	var message string

	if tasks[0].StatusId == 1 {
		message = "Ваши задачи:\n\n"
	} else {
		message = "Ваши завершенные задачи:\n\n"
	}

	for i, task := range tasks {
		if task.StatusId == 1 {
			message += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %s\n\n", i+1, task.Title, task.Description, "в процессе")
		}
		if task.StatusId == 2 {
			message += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %s\n\n", i+1, task.Title, task.Description, "выполнено")
		}
	}

	return &chatID, message
}
