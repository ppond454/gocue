package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type IQueue[T any] interface {
	Process(callback func(data T, ack func()))
	Push(data T) error
	Shutdown()
}

func New[T any](name string, enableAck bool, client *redis.Client) IQueue[T] {
	ctx := context.Background()
	return &RedisQueue[T]{
		Key:       Key(name),
		client:    client,
		ctx:       ctx,
		enableAck: enableAck,
	}
}
