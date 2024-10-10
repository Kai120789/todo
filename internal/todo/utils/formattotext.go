package utils

import (
	"fmt"
	"todo/internal/todo/models"
)

func FormatTasksMessage(tasks []models.Task) string {
	if len(tasks) == 0 {
		return "У вас нет задач."
	}

	message := "Ваши задачи:\n\n"
	for i, task := range tasks {
		message += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %s\n\n", i+1, task.Title, task.Description, "в процессе")
	}

	return message
}

func FormatEndedTasksMessage(tasks []models.Task) string {
	if len(tasks) == 0 {
		return "У вас нет задач, завершенных сегодня."
	}

	message := "Ваши завершенные задачи:\n\n"
	for i, task := range tasks {
		message += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %s\n\n", i+1, task.Title, task.Description, "выполнено")
	}

	return message
}
