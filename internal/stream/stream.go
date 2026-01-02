package stream

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStream struct {
	Client *redis.Client
	Name   string
}

type Stream interface {
	xAdd(ctx context.Context, url string, id int64) error
	XAddBulk(ctx context.Context, websites []WebsiteEvent) error
	XReadGroup(ctx context.Context, consumerGroup, workerID string) ([]redis.XMessage, error)
	XAckBulk(ctx context.Context, consumerGroup string, eventIDs []string) error
}

type WebsiteEvent struct {
	ID int64 `json:"id"`
	Url string `json:"url"`
}

func NewRedisStream(client *redis.Client, name string) *RedisStream {
	return &RedisStream{
		Client: client,
		Name:   name,
	}
}

func OpenRedisConnection(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return client, nil
}

func (s *RedisStream) xAdd(ctx context.Context, url string, id int64) error {
	_, err := s.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: s.Name,
		Values: map[string]interface{}{
			"id": id,
			"url": url,
		},
	}).Result()
	
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStream) XAddBulk(ctx context.Context, websites []*WebsiteEvent) error {
	for _, website := range websites {
		err := s.xAdd(ctx, website.Url, website.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RedisStream) XReadGroup(ctx context.Context, consumerGroup, workerID string) ([]redis.XMessage, error) {
	res, err := s.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: workerID,
		Streams:  []string{s.Name, ">"},
		Count:    5,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0].Messages, nil
}

func (s *RedisStream) XAckBulk(ctx context.Context, consumerGroup string, eventIDs []string) error {
	if len(eventIDs) == 0 {
		return nil
	}

	_, err := s.Client.XAck(ctx, s.Name, consumerGroup, eventIDs...).Result()
	return err
}