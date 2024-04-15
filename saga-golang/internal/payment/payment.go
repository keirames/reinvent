package payment

import (
	"fmt"
	"main/internal/broker"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

func Create() {
	go broker.KafkaCreateReader("order_created", func(id string) {
		fmt.Println("executor got order", id)

		// payment check cost 5 sec
		time.Sleep(time.Second * 5)

		err := broker.KafkaWriteMessages("order_paid", kafka.Message{
			Value: []byte(id),
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("paid event sent", id)
		}
	})

	r := mux.NewRouter()
	url := "0.0.0.0:3000"
	err := http.ListenAndServe(url, r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listen on %v", url)
}
