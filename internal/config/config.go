package config

import (
	"errors"
	"os"
)

type Config struct {
    DB string
}

func Load() (Config, error) {
    v := os.Getenv("DATABASE_URL")
    if v == "" {
        return Config{}, errors.New("DATABASE_URL is required")
    }

    return Config{DB: v}, nil
}
