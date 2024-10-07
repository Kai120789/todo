package tgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db *pgxpool.Pool
}

type Storager interface {
	RegisterUser(upd int, tgName string, chatID int64) error
	SendTask()
	DailyReport()
}

func New(Conn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{db: Conn}
}

func Connection(connectionStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}

func (d *Storage) RegisterUser(upd int, tgName string, chatID int64) error {
	query := `INSERT INTO tg_id (tg_name, chat_id) VALUES ($1, $2)`
	_, err := d.db.Exec(context.Background(), query, tgName, chatID)
	if err != nil {
		fmt.Printf("Error registering user: %v\n", err)
		return err
	}

	return nil
}

func (d *Storage) GetMyTasks() {

}

func CreateTask() {

}

func DailyReport() {

}
