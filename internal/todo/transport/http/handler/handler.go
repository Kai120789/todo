package handler

import (
	"go.uber.org/zap"
)

type TodoHandler struct {
	BoardsHandler   BoardsHandler
	StatusesHandler StatusesHandler
	TasksHandler    TasksHandler
	UserHandler     UserHandler
}

type TodoService struct {
	BoardsService   BoardsHandlerer
	StatusesService StatusesHandlerer
	TasksService    TasksHandlerer
	UserService     UserHandlerer
}

func New(t TodoService, logger *zap.Logger) TodoHandler {
	return TodoHandler{
		BoardsHandler:   NewBoardsHandler(t.BoardsService, logger),
		StatusesHandler: NewStatusesHandler(t.StatusesService, logger),
		TasksHandler:    NewTasksHandler(t.TasksService, logger),
		UserHandler:     NewUserHandler(t.UserService, logger),
	}
}
