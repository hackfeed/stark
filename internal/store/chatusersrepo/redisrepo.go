package chatusersrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hackfeed/stark/internal/domain"
	"github.com/hackfeed/stark/internal/store"
	jsoniter "github.com/json-iterator/go"
)

const (
	prefix = "users"
)

type redisRepo struct {
	ttl         time.Duration
	cacheClient domain.Cacher
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewRedisRepo(cacheClient domain.Cacher, ttl time.Duration) store.UsersRepository {
	return &redisRepo{
		cacheClient: cacheClient,
		ttl:         ttl,
	}
}

func (rr *redisRepo) GetUsers(chat string) (map[string]struct{}, error) {
	key := prefix + chat

	bytes, err := rr.cacheClient.Get(context.Background(), key)
	if err == redis.Nil {
		return nil, nil
	}
	if err != redis.Nil && err != nil {
		return nil, fmt.Errorf("failed to get users by key '%s', error is: %s", key, err)
	}

	users := make(map[string]struct{})

	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes to users on key '%s', error is: %s", key, err)
	}

	return users, nil
}

func (rr *redisRepo) SetUsers(chat string, users map[string]struct{}) error {
	key := prefix + chat

	bytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users to JSON with key '%s', error is: %s", key, err)
	}

	if err := rr.cacheClient.Set(context.Background(), key, bytes, rr.ttl); err != nil {
		return fmt.Errorf("failed to set users with key '%s', error is: %s", key, err)
	}

	return nil
}
