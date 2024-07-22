package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
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
		log.Printf("Received a message: %s", d.Body)
	}
}
