package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://user:pass@localhost:5672/")
	if err != nil {
		panic(err)
	}
	fmt.Println("rabbitmq connected")

	return conn
}

func connectRedis() *redis.Client {
	opt, err := redis.ParseURL("redis://localhost:6379/jobs")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	fmt.Println("redis connected")

	return client
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	conn := connect()
	defer conn.Close()

	redisConn := connectRedis()
	defer redisConn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
