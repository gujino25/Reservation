package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return Config{}, errors.New("DATABASE_URL is not set: copy .env.example to .env")
	}
	return Config{
		DatabaseURL: url,
		HTTPAddr:    envOr("HTTP_ADDR", ":8080"),
	}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
