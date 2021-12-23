package usersrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hackfeed/stark/internal/db/cache"
)

const keyPrefix = "users"

type redisRepo struct {
	ttl         time.Duration
	cacheClient cache.RedisClient
}

func NewRedisRepo(cacheClient cache.RedisClient, ttl time.Duration) UsersRepository {
	return &redisRepo{
		cacheClient: cacheClient,
		ttl:         ttl,
	}
}

func (rr *redisRepo) GetUsers(chat string) (map[string]struct{}, error) {
	bytes, err := rr.cacheClient.Get(context.Background(), fmt.Sprintf("%s:%s", keyPrefix, chat))
	if err == redis.Nil {
		return nil, nil
	}
	if err != redis.Nil && err != nil {
		return nil, fmt.Errorf("failed to get users by key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	users := make(map[string]struct{})

	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes to users on key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	return users, nil
}

func (rr *redisRepo) SetUsers(chat string, users map[string]struct{}) error {
	bytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users to JSON with key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	if err := rr.cacheClient.Set(context.Background(), fmt.Sprintf("%s:%s", keyPrefix, chat), bytes, rr.ttl); err != nil {
		return fmt.Errorf("failed to set users with key '%s', error is: %s", fmt.Sprintf("%s:%s", keyPrefix, chat), err)
	}

	return nil
}
