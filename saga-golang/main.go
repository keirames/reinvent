package main

import (
	"context"
	"fmt"
	"main/internal/broker"
	"main/internal/order"
	"main/internal/payment"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func connectDB() (*pgx.Conn, error) {
	databaseUrl := fmt.Sprintf(
		"postgres://postgres:password@%v:5432/postgres",
		os.Getenv("DATABASE_HOST"),
	)
	for i := 0; i < 10; i++ {
		conn, err := pgx.Connect(context.Background(), databaseUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		} else {
			return conn, nil
		}

		fmt.Println("Retrying...")
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("Database not response")
}

func retryDial() (*kafka.Conn, error) {
	for i := 0; i < 30; i++ {
		conn, err := kafka.Dial("tcp", broker.MapPort(9092))
		if err != nil {
			fmt.Println(err)
		} else {
			return conn, nil
		}

		fmt.Println("Retrying connect to broker...")
		time.Sleep(time.Second * 5)
	}

	return nil, fmt.Errorf("Fail to connect to broker")
}

func prepareTopic() {
	conn, err := retryDial()
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial(
		"tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	)
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             "order_created",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             "order_paid",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		{
			Topic:             "order_delivered",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	conn, err := connectDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Println("database connected!")

	_, err = conn.Exec(
		context.Background(),
		`create table if not exists orders (
			id serial primary key,
			state text check (state in ('created', 'paid', 'delivered'))
		)`,
	)
	if err != nil {
		fmt.Println("Cannot prepare database table")
		os.Exit(1)
	}
	fmt.Println("database prepared!")

	prepareTopic()
	fmt.Println("topic created!")

	state := os.Getenv("STATE")
	if state == "order" {
		order.Create(conn)
	}

	if state == "payment" {
		payment.Create()
	}
}
