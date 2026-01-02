package config

import (
	"errors"
	"os"
)

type Config struct {
    DB string
    RedisAddr string
    RegionID string
    WorkerID string
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

    regionID := os.Getenv("REGION_ID")
    if regionID == "" {
        return Config{}, errors.New("REGION_ID is required")
    }

    workerID := os.Getenv("WORKER_ID")
    if workerID == "" {
        return Config{}, errors.New("WORKER_ID is required")
    }

    return Config{DB: v, RedisAddr: redisAddr, RegionID: regionID, WorkerID: workerID}, nil
}