package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"
	"wb-tech/internal/common/config"
	"wb-tech/internal/pkg/repo"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func Connect(cfg config.Database) (*Storage, error) {
	const op = "Storage.postgre.New"

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		repo.GetEnv("DB_HOST", cfg.Host),
		repo.GetEnv("DB_PORT", cfg.Port),
		repo.GetEnv("DB_USER", cfg.User),
		repo.GetEnv("DB_PASSWORD", cfg.Password),
		repo.GetEnv("DB_NAME", cfg.Dbname),
	)

	log.Printf("Attempting to connect to database with: %s", repo.HidePassword(sqlInfo))

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	migrationsDir, err := filepath.Abs("./migrations")
	// migrationsDir, err := filepath.Abs("../../migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get migrations path: %w", op, err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return nil, fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	}

	return &Storage{db: db}, nil
}
