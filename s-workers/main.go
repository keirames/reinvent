package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func main() {
	// somehow get msg to process

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.Ping(context.Background())
}
