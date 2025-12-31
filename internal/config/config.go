package config

import (
	"errors"
	"os"
)

type Config struct {
    DB string
    RedisAddr string
}

func Load() (Config, error) {
    v := os.Getenv("DATABASE_URL")
    if v == "" {
        return Config{}, errors.New("DATABASE_URL is required")
    }

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        return Config{}, errors.New("REDIS_ADDR is required")
    }

    return Config{DB: v, RedisAddr: redisAddr}, nil
}
