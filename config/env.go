package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "production" {
		return nil
	}

	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}