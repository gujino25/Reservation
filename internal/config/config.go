package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return Config{}, errors.New("DATABASE_URL is not set: copy .env.example to .env")
	}
	return Config{DatabaseURL: url}, nil
}
