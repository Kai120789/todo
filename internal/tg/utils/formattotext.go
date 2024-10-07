package utils

import (
	"fmt"
	"todo/internal/todo/models"
)

func FormatTasksMessage(tasks []models.Task) string {
	if len(tasks) == 0 {
		return "У вас нет задач."
	}

	message := "Ваши задачи:\n"
	for i, task := range tasks {
		message += fmt.Sprintf("%d. %s\nОписание: %s\nСтатус: %d\n\n", i+1, task.Title, task.Description, task.StatusId)
	}

	return message
}
