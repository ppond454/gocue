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

	for i := range 20000 {
		go func() {
			data := TData{
				Id:   i,
				Name: fmt.Sprintf("%d_name", i),
			}
			test.Push(data)
		}()
	}

	test.Shutdown()
}
