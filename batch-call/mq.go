package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
}

type Message struct {
	Type string
}

func New() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://admin:123@localhost:5672/")

	if err != nil {
		return nil, err
	}

	r := new(RabbitMQ)
	r.conn = conn

	return r, nil
}

func (r *RabbitMQ) Consume() {
	ch, err := r.conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	for d := range msgs {
		var m Message
		err := json.Unmarshal(d.Body, &m)
		if err != nil {
			fmt.Println("unknown message, skip")
			continue
		}

		executor, err := sync.GetExecutorByToken(m.Type)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = executor.Execute()
		if err != nil {
			//! don't commit
			continue
		}

		log.Printf("Received a message: %s", d.Body)
	}
}
