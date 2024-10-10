package services

import (
	"go.uber.org/zap"
)

type TodoService struct {
	BoardsService   BoardsService
	StatusesService StatusesService
	TasksService    TasksService
	UserService     UserService
}

type Storager struct {
	BoardsStorager   BoardsStorager
	StatusesStorager StatusesStorager
	TasksStorager    TasksStorager
	UserStorager     UserStorager
}

func New(stor Storager, log *zap.Logger) *TodoService {
	return &TodoService{
		BoardsService:   *NewBoardsService(stor.BoardsStorager, log),
		StatusesService: *NewStatusesService(stor.StatusesStorager, log),
		TasksService:    *NewTasksService(stor.TasksStorager, log),
		UserService:     *NewUserService(stor.UserStorager, log),
	}
}
