package broker

import (
	"context"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

func MapPort(p int) string {
	host := os.Getenv("BROKER_HOST")
	return fmt.Sprintf("%v:%v", host, p)
}

func KafkaWriteMessages(topic string, msg kafka.Message) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(MapPort(9092), MapPort(9093), MapPort(9094)),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	err := w.WriteMessages(
		context.Background(),
		msg,
	)
	if err != nil {
		fmt.Println("failed to write messages:", err)
		return err
	}

	if err := w.Close(); err != nil {
		fmt.Println("failed to close writer:", err)
		return err
	}

	return nil
}

func KafkaCreateReader(topic string, executor func(id string)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{MapPort(9092), MapPort(9093), MapPort(9094)},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		executor(string(m.Value))
	}

	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}
