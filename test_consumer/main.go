package main

import (
	"fmt"

	"github.com/ppond454/gocue/queue"
	"github.com/redis/go-redis/v9"
)

type TData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Change if Redis is remote
		Password: "",               // Set if Redis has a password
		DB:       0,                // Use default DB
	})

	test := queue.New[TData]("test", true, client)

	count := 0
	test.Process(func(data TData, ack func()) {
		count++
		fmt.Println("data: ", data, count)
		ack()
	})
	test.Shutdown()
}
