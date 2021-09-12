package chatmessagerepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hackfeed/stark/internal/domain"
	"github.com/hackfeed/stark/internal/domain/chatmessage"
	"github.com/hackfeed/stark/internal/store"
)

const (
	keyPostfix = "chatmessages"
)

type redisRepo struct {
	ttl         time.Duration
	cacheClient domain.Cacher
}

func NewRedisRepo(cacheClient domain.Cacher, ttl time.Duration) store.MessagesRepository {
	return &redisRepo{
		cacheClient: cacheClient,
		ttl:         ttl,
	}
}

func (rr *redisRepo) GetMessages(id store.Identifier) ([]domain.Messager, error) {
	key := id.FormatIDWithPostfix(keyPostfix)

	bytes, err := rr.cacheClient.Get(context.Background(), key)
	if err == redis.Nil {
		return nil, nil
	}
	if err != redis.Nil && err != nil {
		return nil, fmt.Errorf("failed to get messages by key '%s'", key)
	}

	intMessages := chatmessage.ChatMessages{}

	if err := json.Unmarshal(bytes, &intMessages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes to messages on key '%s', error is: %s", key, err)
	}

	messages := make([]domain.Messager, len(intMessages))
	for i, msg := range intMessages {
		messages[i] = &msg
	}

	return messages, nil
}

func (rr *redisRepo) SetMessages(id store.Identifier, messages []domain.Messager) error {
	key := id.FormatIDWithPostfix(keyPostfix)

	bytes, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal messages to JSON with key '%s'", key)
	}

	if err := rr.cacheClient.Set(context.Background(), key, bytes, rr.ttl); err != nil {
		return fmt.Errorf("failed to set messages with key '%s'", key)
	}

	return nil

}
