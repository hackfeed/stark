package domain

import (
	"context"
	"time"
)

type Cacher interface {
	Set(context.Context, string, interface{}, time.Duration) error
	Get(context.Context, string) ([]byte, error)
	Publish(context.Context, string, interface{}) error
	Subscribe(context.Context, string) <-chan string
}
