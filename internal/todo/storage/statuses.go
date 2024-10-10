package storage

import (
	"context"
	"todo/internal/todo/dto"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type StatusesStorage struct {
	db *pgxpool.Pool
}

type StatusesStorager interface {
	SetStatus(body dto.PostStatusDto) error
	DeleteStatus(id uint) error
}

func NewStatusesStore(Conn *pgxpool.Pool, log *zap.Logger) *StatusesStorage {
	return &StatusesStorage{db: Conn}
}

// create status
func (d *StatusesStorage) SetStatus(body dto.PostStatusDto) error {
	query := `INSERT INTO statuses type VALUES $1`
	_, err := d.db.Exec(context.Background(), query, body.Type)
	if err != nil {
		return err
	}

	return nil
}

// delete status
func (d *StatusesStorage) DeleteStatus(id uint) error {
	query := `DELETE FROM statuses WHERE id = $1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
