package main

import (
	"wb-tech/internal/common/config"
	"wb-tech/internal/pkg/logger"
	"wb-tech/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	logger.Init()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Log.Panic("No .env file found")
	}

	cfg := config.MustLoad()

	_, err := repository.Connect(cfg.Database)
	if err != nil {
		logger.Log.Errorf("failed to init storage: %v", err)
	}

}
