package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func exponentialBackoff[V *amqp.Connection | *redis.Client](
	cb func() (V, error),
) (V, error) {
	waitSec := 0
	for {
		v, err := cb()
		if err != nil {
			return v, nil
		}
		if waitSec > 10 {
			return nil, err
		}

		waitSec += 2
		time.Sleep(time.Second * time.Duration(waitSec))
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	conn, err := exponentialBackoff(func() (*amqp.Connection, error) {
		return amqp.Dial("amqp://user:pass@localhost:5672/")
	})
	if err != nil {
		fmt.Println("rabbitmq connection fail")
		return nil, err
	}

	fmt.Println("rabbitmq connected")
	return conn, nil
}

func connectRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL("redis://localhost:6379/jobs")
	if err != nil {
		fmt.Println("fail to parse redis options")
		return nil, err
	}

	client := redis.NewClient(opt)

	return client, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	rmqConn, err := connectRabbitMQ()
	if err != nil {
		panic(err)
	}
	defer rmqConn.Close()

	redisConn, err := connectRedis()
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()

	ch, err := rmqConn.Channel()
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
