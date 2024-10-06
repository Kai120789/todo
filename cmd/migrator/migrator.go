package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	dbDSN := os.Getenv("DBDSN")

	migrateDsn := dbDSN[:27] + "localhost:5431/taskdb?sslmode=disable"

	fmt.Println(migrateDsn)

	db, err := pgxpool.New(context.Background(), migrateDsn)
	if err != nil {
		zap.S().Fatal("connect db error: ", err)
	}

	var direction string
	flag.StringVar(&direction, "d", "", "direction of migration: 'down' or 'up'")
	flag.Parse()

	if direction == "" {
		err = MigrateUp(db, "./migrations", "up")
		if err != nil {
			return
		}
	} else if direction == "down" {
		err = MigrateDown(db, "./migrations", "down")
		if err != nil {
			return
		}
	}

	fmt.Println("Миграции успешно применены!")
}

func MigrateUp(db *pgxpool.Pool, migrationPath string, direction string) error {
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "up.sql") {
			sqlFilePath := filepath.Join(migrationPath, file.Name())
			err := executeMigration(db, sqlFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func MigrateDown(db *pgxpool.Pool, migrationPath string, direction string) error {
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".down.sql") {
			sqlFilePath := filepath.Join(migrationPath, file.Name())
			err := executeMigration(db, sqlFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func executeMigration(db *pgxpool.Pool, sqlFilePath string) error {
	schemaSQL, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Executing migration: %s\n", sqlFilePath)

	_, err = db.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		fmt.Printf("Ошибка при выполнении миграции %s: %v\n", sqlFilePath, err) // Логирование ошибок
		return err
	}

	return nil
}
