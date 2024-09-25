package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var dbpool *pgxpool.Pool

func InitDB(password, dbname, host string, port int) error {
	databaseUrl := fmt.Sprintf("postgres://postgres:%s@%s:%d/%s", password, host, port, dbname)

	//databaseUrl := "postgres://postgres:yourpassword@0.0.0.0:5432/taskdb"

	ctx := context.Background()

	var err error
	dbpool, err = pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Проверяем подключение
	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	return err
}
