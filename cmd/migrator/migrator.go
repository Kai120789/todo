package migrator

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	dbDSN := os.Getenv("DSN")

	db, err := pgxpool.New(context.Background(), dbDSN)
	if err != nil {
		zap.S().Fatal("connect db error: ", err)
	}

	_ = db

}
