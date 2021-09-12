package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	extredis "github.com/go-redis/redis/v8"
)

var (
	redisClient *RedisClient
	lock        = &sync.Mutex{}
)

type RedisClient struct {
	client *extredis.Client
}

type Options struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient(ctx context.Context, options *Options) (*RedisClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if redisClient == nil {
		client, err := getRedisClient(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize redis client, error is: %s", err)
		}
		redisClient = &RedisClient{
			client: client,
		}
		return redisClient, nil
	}

	return redisClient, nil
}

func getRedisClient(ctx context.Context, options *Options) (*extredis.Client, error) {
	opts := extredis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	}
	client := extredis.NewClient(&opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis client with ttl %s failed to ping address %s, error is: %s",
			opts.IdleTimeout, opts.Addr, err)
	}

	return client, nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return rc.client.Get(ctx, key).Bytes()
}

func (rc *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return rc.client.Publish(ctx, channel, message).Err()
}

func (rc *RedisClient) Subscribe(ctx context.Context, channel string) <-chan string {
	redisChan := rc.client.Subscribe(ctx, channel).Channel()
	messagesChan := make(chan string)
	go func(ch <-chan *extredis.Message) {
		for msg := range ch {
			messagesChan <- msg.Payload
		}
	}(redisChan)

	return messagesChan
}
