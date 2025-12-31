package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adi-ty/betteruptime-go/internal/config"
	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/Adi-ty/betteruptime-go/internal/stream"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    db, err := store.Open(cfg.DB)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    client, err := stream.OpenRedisConnection(cfg.RedisAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    streamInstance := stream.NewRedisStream(client, "Betteruptime:Websites")

    ticker := time.NewTicker(3 * time.Minute)
    defer ticker.Stop()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    log.Println("Pusher started. Press Ctrl+C to stop.")

    for {
        select {
        case <-ticker.C:
            pushWebsites(db, streamInstance)
        case <-sigChan:
            log.Println("Received interrupt signal. Shutting down...")
            return
        }
    }
}

func pushWebsites(db *sql.DB, streamInstance *stream.RedisStream) {
    websiteStore := store.NewPostgresWebsiteStore(db)

    websites, err := websiteStore.GetAllWebsites()
    if err != nil {
        log.Printf("Error fetching websites: %v", err)
        return
    }

    err = streamInstance.XAddBulk(context.Background(), websites)
    if err != nil {
        log.Printf("Error pushing to Redis: %v", err)
        return
    }

    log.Println("Websites pushed to Redis stream")
}