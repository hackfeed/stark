package filesrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hackfeed/stark/internal/db/cache"
)

const keyPrefix = "files"

type redisRepo struct {
	ttl         time.Duration
	cacheClient cache.RedisClient
}

func NewRedisRepo(cacheClient cache.RedisClient, ttl time.Duration) FilesRepository {
	return &redisRepo{
		cacheClient: cacheClient,
		ttl:         ttl,
	}
}

func (rr *redisRepo) GetFiles(chat string) (map[string][]byte, error) {
	bytes, err := rr.cacheClient.Get(context.Background(), fmt.Sprintf("%s:%s", keyPrefix, chat))
	if err == redis.Nil {
		return nil, nil
	}
	if err != redis.Nil && err != nil {
		return nil, fmt.Errorf("failed to get files by key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	files := make(map[string][]byte)

	if err := json.Unmarshal(bytes, &files); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes to files on key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	return files, nil
}

func (rr *redisRepo) SetFiles(chat string, files map[string][]byte) error {
	bytes, err := json.Marshal(files)
	if err != nil {
		return fmt.Errorf("failed to marshal files to JSON with key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	if err := rr.cacheClient.Set(context.Background(), fmt.Sprintf("%s:%s", keyPrefix, chat), bytes, rr.ttl); err != nil {
		return fmt.Errorf("failed to set files with key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	return nil
}
