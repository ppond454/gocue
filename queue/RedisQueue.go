package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

type Key string

func (k Key) KeyPath() string {
	return fmt.Sprintf("go_cue:%s:queue", k)
}

type RedisQueue[T any] struct {
	Key       Key `json:"key"`
	client    *redis.Client
	ctx       context.Context
	enableAck bool
}

func (q *RedisQueue[T]) Process(callback func(data T, ack func())) {

	thread := func() {
		for {
			task, err := q.client.BLPop(q.ctx, 0, q.Key.KeyPath()).Result()
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			raw := []byte(task[1])
			item, errDecode := q.decode(raw)

			if errDecode != nil {
				fmt.Println("Error:", errDecode)
				return
			}
			isAck := false
			callback(item, func() {
				isAck = true
			})
			if q.enableAck && !isAck {
				q.rollBack(item)
			}
		}
	}

	go thread()
}

func (q *RedisQueue[T]) encode(value T) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (q *RedisQueue[T]) decode(value []byte) (T, error) {
	var data T
	err := json.Unmarshal(value, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (q *RedisQueue[T]) rollBack(data T) error {
	item, _ := q.encode(data)
	log.Println("rollback", item)
	return q.client.LPush(q.ctx, q.Key.KeyPath(), item).Err()
}

func (q *RedisQueue[T]) Push(data T) error {
	_data, err := q.encode(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return q.client.RPush(q.ctx, q.Key.KeyPath(), _data).Err()
}

func (q *RedisQueue[T]) Shutdown() {
	ctx, cancel := context.WithCancel(q.ctx)
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Program is shutting down gracefully...")
				return
			default:
				return
			}
		}
	}()

	sigReceived := <-signalChan
	fmt.Printf("Received signal: %s\n", sigReceived)

	cancel()
	fmt.Println("Program exited cleanly.")
}
