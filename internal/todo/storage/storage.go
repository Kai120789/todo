package storage

import (
	"context"
	"fmt"

	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(Conn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{db: Conn, logger: log}
}

func Connection(connectionStr string, logger *zap.Logger) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	log.Println("Successfully connect to db")

	return db, nil
}

func (d *Storage) Close() error {
	if d.db == nil {
		return nil
	}
	d.db.Close()
	return nil
}
