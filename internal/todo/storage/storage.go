package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	BoardsStorage   BoardsStorage
	TasksStorage    TasksStorage
	StatusesStorage StatusesStorage
	UserStorage     UserStorage
}

func New(Conn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{
		BoardsStorage:   *NewBoardsStore(Conn, log),
		TasksStorage:    *NewTasksStore(Conn, log),
		StatusesStorage: *NewStatusesStore(Conn, log),
		UserStorage:     *NewUserStore(Conn, log),
	}
}

func Connection(connectionStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}
