package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Adi-ty/betteruptime-go/internal/config"
	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/Adi-ty/betteruptime-go/internal/stream"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type fetchResult struct {
    eventID string
    success bool
}

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

	
	storeInstance := store.NewPostgresWebsiteStore(db)
	streamInstance := stream.NewRedisStream(client, "Betteruptime:Websites")

	for {
		messages, err := streamInstance.XReadGroup(context.Background(), cfg.RegionID, cfg.WorkerID)
		if err == redis.Nil {
			continue
		}
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		if len(messages) == 0 {
			continue
		}

		results := make(chan fetchResult, len(messages))

		started := 0
		for _, msg := range messages {
			url, ok := msg.Values["url"].(string)
			if !ok {
				log.Printf("invalid url in message %s: %v", msg.ID, msg.Values["url"])
				continue
			}
			idStr, ok := msg.Values["id"].(string)
			if !ok {
				log.Printf("invalid id in message %s: %v", msg.ID, msg.Values["id"])
				continue
			}
			websiteID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Printf("failed to parse websiteID in message %s: %v", msg.ID, err)
				continue
			}

			started++
			go func(url string, websiteID int64, eventID string) {
				defer func() {
					if r := recover(); r != nil {
                        log.Printf("goroutine panic for website %d: %v", websiteID, r)
                        results <- fetchResult{eventID: eventID, success: false}
                    }
				}()
				success := fetchWebsite(storeInstance, url, cfg.RegionID, websiteID)
				results <- fetchResult{eventID: eventID, success: success}
			}(url, websiteID, msg.ID)
		}

		eventIDs := make([]string, 0, started)
		for i := 0; i < started; i++ {
			result := <-results
			if result.success {
				eventIDs = append(eventIDs, result.eventID)
			}
		}

		log.Println("processed", len(eventIDs), "successful out of", started, "attempted")

		if err := streamInstance.XAckBulk(context.Background(), cfg.RegionID, eventIDs); err != nil {
			log.Println("ack error:", err)
		}
	}
}

func fetchWebsite(storeInstance store.WebsiteStore, url string, RegionID string, websiteID int64) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	startTime := time.Now()
	resp, err := client.Get(url)
	responseTime := time.Since(startTime)

	var websiteTick store.WebsiteTick
	websiteTick.ID = uuid.New().String()
	websiteTick.WebsiteID = websiteID
	websiteTick.RegionID = RegionID
	websiteTick.ResponseTimeMs = responseTime.Milliseconds()

	if err != nil {
		log.Printf("error fetching website %d (%s): %v", websiteID, url, err)
		websiteTick.StatusCode = store.StatusUnknown
	} else {
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			websiteTick.StatusCode = store.StatusUp
		} else {
			websiteTick.StatusCode = store.StatusDown
		}
		log.Printf("fetched website %d (%s): status %s, time %v", websiteID, url, websiteTick.StatusCode, responseTime)
	}

	if err := storeInstance.MarkWebsiteTickProcessed(&websiteTick); err != nil {
		log.Printf("error marking website tick processed for %d: %v", websiteID, err)
		return false
	}

	return true
}